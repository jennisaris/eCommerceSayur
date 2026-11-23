package usecase

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"
	"testing"
)

type fakeOrderRepoForAccess struct {
	orders map[string]*domain.Order
}

func (r *fakeOrderRepoForAccess) CreateWithItems(_ context.Context, _ *domain.Order, _ []domain.OrderItem) error {
	return nil
}

func (r *fakeOrderRepoForAccess) GetByID(_ context.Context, _ int64) (*domain.Order, error) {
	return nil, domain.ErrNotFound
}

func (r *fakeOrderRepoForAccess) GetByCode(_ context.Context, code string) (*domain.Order, error) {
	o, ok := r.orders[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (r *fakeOrderRepoForAccess) GetByCodeWithDetails(_ context.Context, code string) (*domain.Order, error) {
	o, ok := r.orders[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return o, nil
}

func (r *fakeOrderRepoForAccess) GetAll(_ context.Context, _ domain.OrderFilter) ([]domain.Order, int64, error) {
	return nil, 0, nil
}

func (r *fakeOrderRepoForAccess) UpdateStatus(_ context.Context, _ int64, _ string, _ string) error {
	return nil
}

func (r *fakeOrderRepoForAccess) RestockCancelledOrder(_ context.Context, _ int64) error {
	return nil
}

func TestGetOrderByCodeOwnershipAndFreshAdminRole(t *testing.T) {
	owner := &domain.User{
		ID:    100,
		Email: "owner@example.com",
		Roles: []domain.Role{{ID: 2, Name: "Customer"}},
	}
	otherCustomer := &domain.User{
		ID:    200,
		Email: "other@example.com",
		Roles: []domain.Role{{ID: 2, Name: "Customer"}},
	}
	adminUser := &domain.User{
		ID:    300,
		Email: "admin@example.com",
		Roles: []domain.Role{{ID: 1, Name: "Super Admin"}},
	}

	userRepo := newFakeAuthUserRepo(owner)
	userRepo.users[otherCustomer.ID] = otherCustomer
	userRepo.users[adminUser.ID] = adminUser

	orderRepo := &fakeOrderRepoForAccess{
		orders: map[string]*domain.Order{
			"INV-20261123-AAA111": {
				ID:        1,
				OrderCode: "INV-20261123-AAA111",
				BuyerID:   owner.ID,
				Status:    string(domain.OrderStatusPendingPayment),
			},
		},
	}

	orderUC := NewOrderUseCase(orderRepo, nil, nil, userRepo, nil, nil)

	// 1. Owner can read own order
	got, err := orderUC.GetOrderByCode(context.Background(), owner.ID, "INV-20261123-AAA111")
	if err != nil || got.ID != 1 {
		t.Fatalf("expected owner to read own order, got order=%v err=%v", got, err)
	}

	// 2. Other customer is blocked (IDOR prevention)
	if _, err := orderUC.GetOrderByCode(context.Background(), otherCustomer.ID, "INV-20261123-AAA111"); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden for non-owner customer, got %v", err)
	}

	// 3. Active Super Admin in DB can read the order
	gotAdmin, err := orderUC.GetOrderByCode(context.Background(), adminUser.ID, "INV-20261123-AAA111")
	if err != nil || gotAdmin.ID != 1 {
		t.Fatalf("expected Super Admin to read order, got order=%v err=%v", gotAdmin, err)
	}

	// 4. Revoking Super Admin role in DB immediately blocks access
	adminUser.Roles = []domain.Role{{ID: 2, Name: "Customer"}}
	if _, err := orderUC.GetOrderByCode(context.Background(), adminUser.ID, "INV-20261123-AAA111"); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden after Super Admin role revocation, got %v", err)
	}
}

func TestCreateSnapTransactionOwnershipCheck(t *testing.T) {
	orderRepo := &fakeOrderRepoForAccess{
		orders: map[string]*domain.Order{
			"INV-20261123-BBB222": {
				ID:        2,
				OrderCode: "INV-20261123-BBB222",
				BuyerID:   100,
				Status:    string(domain.OrderStatusPendingPayment),
			},
		},
	}
	paymentUC := NewPaymentUseCase(nil, orderRepo, nil, nil, nil)

	if _, err := paymentUC.CreateSnapTransaction(context.Background(), "INV-20261123-BBB222", 200); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden when non-owner initiates Snap payment, got %v", err)
	}
}
