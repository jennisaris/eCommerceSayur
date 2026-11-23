package middleware

import (
	"context"
	"ecommerce-backend/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type stubRoleUserRepo struct {
	users map[int64]*domain.User
}

func (r *stubRoleUserRepo) Create(_ context.Context, _ *domain.User) error          { return nil }
func (r *stubRoleUserRepo) GetByEmail(_ context.Context, _ string) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (r *stubRoleUserRepo) Update(_ context.Context, _ *domain.User) error          { return nil }
func (r *stubRoleUserRepo) AssignRole(_ context.Context, _ int64, _ int) error      { return nil }
func (r *stubRoleUserRepo) CreateVerificationToken(_ context.Context, _ *domain.VerificationToken) error {
	return nil
}
func (r *stubRoleUserRepo) GetVerificationToken(_ context.Context, _, _ string) (*domain.VerificationToken, error) {
	return nil, domain.ErrNotFound
}
func (r *stubRoleUserRepo) DeleteVerificationToken(_ context.Context, _ int64) error { return nil }
func (r *stubRoleUserRepo) DeleteVerificationTokensByUser(_ context.Context, _ int64, _ string) error {
	return nil
}
func (r *stubRoleUserRepo) CreateAuthSession(_ context.Context, _ *domain.AuthSession) error {
	return nil
}
func (r *stubRoleUserRepo) GetActiveAuthSession(_ context.Context, _ string) (*domain.AuthSession, error) {
	return nil, domain.ErrInvalidSession
}
func (r *stubRoleUserRepo) RevokeAuthSession(_ context.Context, _ string) error    { return nil }
func (r *stubRoleUserRepo) RevokeAllUserSessions(_ context.Context, _ int64) error { return nil }

func (r *stubRoleUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func TestRequireRoleEnforcesDatabaseRoleAndRevocation(t *testing.T) {
	adminUser := &domain.User{
		ID:    1,
		Email: "admin@example.com",
		Roles: []domain.Role{{ID: 1, Name: "Super Admin"}},
	}
	customerUser := &domain.User{
		ID:    2,
		Email: "customer@example.com",
		Roles: []domain.Role{{ID: 2, Name: "Customer"}},
	}
	repo := &stubRoleUserRepo{
		users: map[int64]*domain.User{
			adminUser.ID:    adminUser,
			customerUser.ID: customerUser,
		},
	}

	buildApp := func(userID int64, setUser bool) *fiber.App {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			if setUser {
				c.Locals("user_id", userID)
				// Even if JWT locals claim Super Admin, RequireRole must verify DB.
				c.Locals("roles", []string{"Super Admin"})
			}
			return c.Next()
		})
		app.Get("/admin/check", RequireRole("Super Admin", repo), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
		return app
	}

	// 1. Unauthenticated -> 401
	resp, err := buildApp(0, false).Test(httptest.NewRequest(http.MethodGet, "/admin/check", nil))
	if err != nil || resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401 for missing user_id, got status=%d err=%v", resp.StatusCode, err)
	}

	// 2. Customer with spoofed/stale JWT role claim -> 403
	resp, err = buildApp(customerUser.ID, true).Test(httptest.NewRequest(http.MethodGet, "/admin/check", nil))
	if err != nil || resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403 for customer, got status=%d err=%v", resp.StatusCode, err)
	}

	// 3. Real Super Admin in DB -> 200
	resp, err = buildApp(adminUser.ID, true).Test(httptest.NewRequest(http.MethodGet, "/admin/check", nil))
	if err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 for Super Admin, got status=%d err=%v", resp.StatusCode, err)
	}

	// 4. After revoking Super Admin role in DB -> 403 immediately
	adminUser.Roles = []domain.Role{{ID: 2, Name: "Customer"}}
	resp, err = buildApp(adminUser.ID, true).Test(httptest.NewRequest(http.MethodGet, "/admin/check", nil))
	if err != nil || resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected 403 after role revocation, got status=%d err=%v", resp.StatusCode, err)
	}
}
