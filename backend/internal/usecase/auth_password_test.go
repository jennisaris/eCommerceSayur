package usecase

import (
	"context"
	"crypto/sha256"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/hash"
	"ecommerce-backend/internal/pkg/jwt"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestHasStrongPassword(t *testing.T) {
	valid := []string{"BookVerse#2026", "安全BookVerse#2026"}
	for _, password := range valid {
		if !hasStrongPassword(password) {
			t.Errorf("expected password %q to satisfy the policy", password)
		}
	}

	invalid := []string{
		"short",
		"alllowercase123!",
		"ALLUPPERCASE123!",
		"NoSpecial1234",
		strings.Repeat("安", 25) + "Aa1!", // > 72 UTF-8 bytes (bcrypt limit)
	}
	for _, password := range invalid {
		if hasStrongPassword(password) {
			t.Errorf("expected password %q to fail the policy", password)
		}
	}
}

type fakeAuthUserRepo struct {
	users          map[int64]*domain.User
	tokens         map[int64]*domain.VerificationToken
	nextTokenID    int64
	revokedUserIDs []int64
}

func newFakeAuthUserRepo(user *domain.User) *fakeAuthUserRepo {
	return &fakeAuthUserRepo{
		users:       map[int64]*domain.User{user.ID: user},
		tokens:      map[int64]*domain.VerificationToken{},
		nextTokenID: 1,
	}
}

func (r *fakeAuthUserRepo) Create(_ context.Context, u *domain.User) error {
	r.users[u.ID] = u
	return nil
}

func (r *fakeAuthUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (r *fakeAuthUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	target := strings.ToLower(strings.TrimSpace(email))
	for _, u := range r.users {
		if strings.ToLower(u.Email) == target {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakeAuthUserRepo) Update(_ context.Context, u *domain.User) error {
	r.users[u.ID] = u
	return nil
}

func (r *fakeAuthUserRepo) AssignRole(_ context.Context, _ int64, _ int) error { return nil }

func (r *fakeAuthUserRepo) CreateVerificationToken(_ context.Context, token *domain.VerificationToken) error {
	token.ID = r.nextTokenID
	r.nextTokenID++
	cp := *token
	r.tokens[cp.ID] = &cp
	return nil
}

func (r *fakeAuthUserRepo) GetVerificationToken(_ context.Context, rawToken string, tokenType string) (*domain.VerificationToken, error) {
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])
	for _, tok := range r.tokens {
		if tok.Token == tokenHash && tok.TokenType == tokenType && tok.ExpiresAt.After(time.Now()) {
			cp := *tok
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (r *fakeAuthUserRepo) DeleteVerificationToken(_ context.Context, id int64) error {
	delete(r.tokens, id)
	return nil
}

func (r *fakeAuthUserRepo) DeleteVerificationTokensByUser(_ context.Context, userID int64, tokenType string) error {
	for id, tok := range r.tokens {
		if tok.UserID == userID && tok.TokenType == tokenType {
			delete(r.tokens, id)
		}
	}
	return nil
}

func (r *fakeAuthUserRepo) CreateAuthSession(_ context.Context, _ *domain.AuthSession) error {
	return nil
}

func (r *fakeAuthUserRepo) GetActiveAuthSession(_ context.Context, _ string) (*domain.AuthSession, error) {
	return nil, domain.ErrInvalidSession
}

func (r *fakeAuthUserRepo) RevokeAuthSession(_ context.Context, _ string) error { return nil }

func (r *fakeAuthUserRepo) RevokeAllUserSessions(_ context.Context, userID int64) error {
	r.revokedUserIDs = append(r.revokedUserIDs, userID)
	return nil
}

type fakeResetEmailSender struct {
	configured bool
	lastURL    string
	sentCount  int
}

func (s *fakeResetEmailSender) IsConfigured() bool { return s.configured }

func (s *fakeResetEmailSender) SendPasswordResetEmail(_ context.Context, _, _, resetURL string) error {
	s.lastURL = resetURL
	s.sentCount++
	return nil
}

func TestPasswordResetFlowAndSessionRevocation(t *testing.T) {
	oldHash, err := hash.HashPassword("OldPassword#2026")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := &domain.User{
		ID:         10,
		Name:       "Budi",
		Email:      "budi@example.com",
		Password:   oldHash,
		IsVerified: true,
	}
	repo := newFakeAuthUserRepo(user)
	emailSender := &fakeResetEmailSender{configured: true}
	uc := NewAuthUseCaseWithOptions(repo, nil, jwt.NewJWTServiceWithDuration("test-secret-32-characters-minimum!", 15*time.Minute), AuthOptions{
		PasswordResetEmailSender: emailSender,
		VerificationBaseURL:      "https://bookverse.example",
		IsProduction:             true,
	})

	// Unknown email must succeed silently (anti user-enumeration)
	if err := uc.RequestPasswordReset(context.Background(), "unknown@example.com"); err != nil {
		t.Fatalf("expected nil error for unknown email, got %v", err)
	}
	if emailSender.sentCount != 0 {
		t.Fatalf("expected no email sent for unknown email, got %d", emailSender.sentCount)
	}

	// Known email (case-insensitive) sends reset link
	if err := uc.RequestPasswordReset(context.Background(), "  BUDI@Example.com "); err != nil {
		t.Fatalf("expected password reset request to succeed, got %v", err)
	}
	if emailSender.sentCount != 1 {
		t.Fatalf("expected 1 reset email sent, got %d", emailSender.sentCount)
	}

	parsed, err := url.Parse(emailSender.lastURL)
	if err != nil {
		t.Fatalf("invalid reset URL %q: %v", emailSender.lastURL, err)
	}
	rawToken := parsed.Query().Get("token")
	if rawToken == "" {
		t.Fatalf("expected non-empty token in reset URL %q", emailSender.lastURL)
	}

	// Reject weak password
	if err := uc.ResetPassword(context.Background(), rawToken, "weak"); !errors.Is(err, domain.ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword, got %v", err)
	}

	// Reset with strong password updates hash, deletes token, and revokes all sessions
	newPass := "NewStrongPass#2026"
	if err := uc.ResetPassword(context.Background(), rawToken, newPass); err != nil {
		t.Fatalf("expected password reset to succeed, got %v", err)
	}
	if !hash.CheckPassword(newPass, repo.users[user.ID].Password) {
		t.Fatal("expected user password hash to match new password")
	}
	if len(repo.revokedUserIDs) != 1 || repo.revokedUserIDs[0] != user.ID {
		t.Fatalf("expected all sessions for user %d to be revoked, got %v", user.ID, repo.revokedUserIDs)
	}

	// Reusing the same token must fail
	if err := uc.ResetPassword(context.Background(), rawToken, "AnotherPass#2026"); !errors.Is(err, domain.ErrInvalidResetToken) {
		t.Fatalf("expected reused token to fail with ErrInvalidResetToken, got %v", err)
	}
}

func TestUpdateProfileAndEmailChangeRequiresCurrentPassword(t *testing.T) {
	pwHash, err := hash.HashPassword("CurrentPass#2026")
	if err != nil {
		t.Fatalf("failed hashing password: %v", err)
	}
	user1 := &domain.User{
		ID:         1,
		Name:       "Budi",
		Email:      "budi@example.com",
		Password:   pwHash,
		Phone:      "081234567890",
		Address:    "Jl. Lama No. 1",
		IsVerified: true,
	}
	user2 := &domain.User{
		ID:         2,
		Name:       "Siti",
		Email:      "siti@example.com",
		Password:   pwHash,
		IsVerified: true,
	}
	repo := newFakeAuthUserRepo(user1)
	repo.users[user2.ID] = user2
	uc := NewAuthUseCase(repo, nil, jwt.NewJWTServiceWithDuration("test-secret-32-characters-minimum!", 15*time.Minute))

	// 1. Update phone & address without changing email succeeds without current_password
	updated, emailChanged, err := uc.UpdateProfile(context.Background(), user1.ID, &domain.UpdateProfileRequest{
		Name:    "Budi Santoso",
		Email:   "budi@example.com",
		Phone:   "081999888777",
		Address: "Jl. Baru No. 99, Jakarta",
	})
	if err != nil || emailChanged {
		t.Fatalf("expected phone/address update to succeed without email change, got err=%v emailChanged=%v", err, emailChanged)
	}
	if updated.Address != "Jl. Baru No. 99, Jakarta" || updated.Phone != "081999888777" {
		t.Fatalf("unexpected updated profile: %+v", updated)
	}

	// 2. Changing email without valid current_password must fail with ErrInvalidCredentials
	if _, _, err := uc.UpdateProfile(context.Background(), user1.ID, &domain.UpdateProfileRequest{
		Name:    "Budi Santoso",
		Email:   "budi.baru@example.com",
		Phone:   "081999888777",
		Address: "Jl. Baru No. 99, Jakarta",
	}); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials when changing email without password, got %v", err)
	}

	// 3. Changing email to another user's email must fail with ErrEmailAlreadyExists
	if _, _, err := uc.UpdateProfile(context.Background(), user1.ID, &domain.UpdateProfileRequest{
		Name:            "Budi Santoso",
		Email:           "SITI@example.com",
		CurrentPassword: "CurrentPass#2026",
	}); !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists for taken email, got %v", err)
	}

	// 4. Changing email with valid current_password succeeds and normalizes email
	updated, emailChanged, err = uc.UpdateProfile(context.Background(), user1.ID, &domain.UpdateProfileRequest{
		Name:            "Budi Santoso",
		Email:           "  Budi.Baru@Example.com ",
		Phone:           "081999888777",
		Address:         "Jl. Baru No. 99, Jakarta",
		CurrentPassword: "CurrentPass#2026",
	})
	if err != nil || !emailChanged || updated.Email != "budi.baru@example.com" {
		t.Fatalf("expected email change to succeed, got user=%+v emailChanged=%v err=%v", updated, emailChanged, err)
	}
}
