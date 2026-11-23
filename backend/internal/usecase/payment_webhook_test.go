package usecase

import (
	"context"
	"crypto/sha512"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/midtrans"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhookEventKeyIsDeterministic(t *testing.T) {
	firstPayload := []byte(`{"order_id":"INV-001","transaction_status":"settlement"}`)
	secondPayload := []byte(`{"order_id":"INV-002","transaction_status":"settlement"}`)

	firstKey := webhookEventKey(firstPayload)
	if firstKey != webhookEventKey(firstPayload) {
		t.Fatal("expected identical webhook payloads to produce the same event key")
	}
	if firstKey == webhookEventKey(secondPayload) {
		t.Fatal("expected different webhook payloads to produce different event keys")
	}
	if len(firstKey) != 64 {
		t.Fatalf("expected SHA-256 hex key length 64, got %d", len(firstKey))
	}
}

type stubPaymentRepoForWebhook struct {
	payments map[int64]*domain.Payment
	logs     map[string]*domain.PaymentLog
}

func (r *stubPaymentRepoForWebhook) Create(_ context.Context, p *domain.Payment) error {
	p.ID = int64(len(r.payments) + 1)
	r.payments[p.OrderID] = p
	return nil
}
func (r *stubPaymentRepoForWebhook) GetByID(_ context.Context, _ int64) (*domain.Payment, error) {
	return nil, domain.ErrNotFound
}
func (r *stubPaymentRepoForWebhook) GetByOrderID(_ context.Context, orderID int64) (*domain.Payment, error) {
	p, ok := r.payments[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return p, nil
}
func (r *stubPaymentRepoForWebhook) GetByGatewayID(_ context.Context, _ string) (*domain.Payment, error) {
	return nil, domain.ErrNotFound
}
func (r *stubPaymentRepoForWebhook) Update(_ context.Context, p *domain.Payment) error {
	r.payments[p.OrderID] = p
	return nil
}
func (r *stubPaymentRepoForWebhook) CreateLog(_ context.Context, l *domain.PaymentLog) error {
	r.logs[l.EventKey] = l
	return nil
}
func (r *stubPaymentRepoForWebhook) GetLogByEventKey(_ context.Context, key string) (*domain.PaymentLog, error) {
	l, ok := r.logs[key]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return l, nil
}

type stubOrderRepoForWebhook struct {
	order         *domain.Order
	updatedStatus string
	restocked     bool
}

func (r *stubOrderRepoForWebhook) CreateWithItems(_ context.Context, _ *domain.Order, _ []domain.OrderItem) error {
	return nil
}
func (r *stubOrderRepoForWebhook) GetByID(_ context.Context, _ int64) (*domain.Order, error) {
	return r.order, nil
}
func (r *stubOrderRepoForWebhook) GetByCode(_ context.Context, code string) (*domain.Order, error) {
	if r.order != nil && r.order.OrderCode == code {
		return r.order, nil
	}
	return nil, domain.ErrNotFound
}
func (r *stubOrderRepoForWebhook) GetByCodeWithDetails(_ context.Context, code string) (*domain.Order, error) {
	return r.GetByCode(context.Background(), code)
}
func (r *stubOrderRepoForWebhook) GetAll(_ context.Context, _ domain.OrderFilter) ([]domain.Order, int64, error) {
	return nil, 0, nil
}
func (r *stubOrderRepoForWebhook) UpdateStatus(_ context.Context, _ int64, status string, _ string) error {
	r.updatedStatus = status
	r.order.Status = status
	return nil
}
func (r *stubOrderRepoForWebhook) RestockCancelledOrder(_ context.Context, _ int64) error {
	r.restocked = true
	r.order.Status = string(domain.OrderStatusCancelled)
	return nil
}

func signMidtransWebhook(orderID, statusCode, grossAmount, serverKey string) string {
	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	return hex.EncodeToString(sum[:])
}

func TestHandleMidtransCallbackVerificationAndStateGuards(t *testing.T) {
	const serverKey = "SB-Mid-server-TestWebhookKey"
	midtransClient := midtrans.NewClient(serverKey, "SB-Mid-client-Test", false)

	buildPayload := func(orderID, status, grossAmount, sig string) []byte {
		body, _ := json.Marshal(domain.MidtransWebhookPayload{
			OrderID:           orderID,
			StatusCode:        "200",
			GrossAmount:       grossAmount,
			SignatureKey:      sig,
			TransactionID:     "tx-123",
			TransactionStatus: status,
			PaymentType:       "qris",
			Currency:          "IDR",
		})
		return body
	}

	// 1. Invalid signature must be rejected
	orderRepo := &stubOrderRepoForWebhook{
		order: &domain.Order{ID: 1, OrderCode: "INV-001", BuyerID: 10, TotalAmount: 100000, Status: string(domain.OrderStatusPendingPayment)},
	}
	payRepo := &stubPaymentRepoForWebhook{payments: map[int64]*domain.Payment{}, logs: map[string]*domain.PaymentLog{}}
	uc := NewPaymentUseCase(payRepo, orderRepo, nil, nil, midtransClient)

	if err := uc.HandleMidtransCallback(context.Background(), buildPayload("INV-001", "settlement", "100000.00", "forged-sig")); !errors.Is(err, domain.ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}

	// 2. Tampered gross amount (even with valid signature for tampered amount) must be rejected
	tamperedSig := signMidtransWebhook("INV-001", "200", "1000.00", serverKey)
	if err := uc.HandleMidtransCallback(context.Background(), buildPayload("INV-001", "settlement", "1000.00", tamperedSig)); err == nil {
		t.Fatal("expected mismatched gross_amount to be rejected")
	}

	// 3. Valid settlement transitions PENDING_PAYMENT -> PAID
	validSig := signMidtransWebhook("INV-001", "200", "100000.00", serverKey)
	if err := uc.HandleMidtransCallback(context.Background(), buildPayload("INV-001", "settlement", "100000.00", validSig)); err != nil {
		t.Fatalf("expected valid settlement webhook to succeed, got %v", err)
	}
	if orderRepo.order.Status != string(domain.OrderStatusPaid) {
		t.Fatalf("expected order status PAID, got %s", orderRepo.order.Status)
	}

	// 4. Late expire webhook on already PAID order must NOT cancel or restock
	if err := uc.HandleMidtransCallback(context.Background(), buildPayload("INV-001", "expire", "100000.00", validSig)); err != nil {
		t.Fatalf("expected late expire webhook to be handled safely, got %v", err)
	}
	if orderRepo.restocked || orderRepo.order.Status != string(domain.OrderStatusPaid) {
		t.Fatalf("late expire must not restock or revert PAID order (restocked=%v, status=%s)", orderRepo.restocked, orderRepo.order.Status)
	}

	// 5. Late settlement webhook on already CANCELLED order must NOT mark it PAID without stock
	cancelledRepo := &stubOrderRepoForWebhook{
		order: &domain.Order{ID: 2, OrderCode: "INV-002", BuyerID: 10, TotalAmount: 50000, Status: string(domain.OrderStatusCancelled)},
	}
	cancelledPayRepo := &stubPaymentRepoForWebhook{payments: map[int64]*domain.Payment{}, logs: map[string]*domain.PaymentLog{}}
	uc2 := NewPaymentUseCase(cancelledPayRepo, cancelledRepo, nil, nil, midtransClient)
	sig2 := signMidtransWebhook("INV-002", "200", "50000.00", serverKey)
	if err := uc2.HandleMidtransCallback(context.Background(), buildPayload("INV-002", "settlement", "50000.00", sig2)); err != nil {
		t.Fatalf("expected late settlement on cancelled order to be logged without crashing, got %v", err)
	}
	if cancelledRepo.order.Status != string(domain.OrderStatusCancelled) {
		t.Fatalf("cancelled order must remain CANCELLED, got %s", cancelledRepo.order.Status)
	}
}

func TestSyncOrderPaymentStatusServerToServer(t *testing.T) {
	const serverKey = "SB-Mid-server-SyncTestKey"
	const orderCode = "INV-20261123-SYNC01"

	// Mock Midtrans Core API GET /v2/{order_id}/status
	mockMidtrans := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/"+orderCode+"/status" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		payload := domain.MidtransWebhookPayload{
			OrderID:           orderCode,
			StatusCode:        "200",
			GrossAmount:       "115000.00",
			SignatureKey:      signMidtransWebhook(orderCode, "200", "115000.00", serverKey),
			TransactionID:     "tx-sync-001",
			TransactionStatus: "settlement",
			PaymentType:       "bank_transfer",
			Currency:          "IDR",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(payload)
	}))
	defer mockMidtrans.Close()

	midtransClient := midtrans.NewClient(serverKey, "SB-Mid-client-SyncTest", false)
	midtransClient.SetAPIBaseURL(mockMidtrans.URL)

	orderRepo := &stubOrderRepoForWebhook{
		order: &domain.Order{
			ID:          10,
			OrderCode:   orderCode,
			BuyerID:     50,
			TotalAmount: 115000,
			Status:      string(domain.OrderStatusPendingPayment),
		},
	}
	payRepo := &stubPaymentRepoForWebhook{payments: map[int64]*domain.Payment{}, logs: map[string]*domain.PaymentLog{}}
	uc := NewPaymentUseCase(payRepo, orderRepo, nil, nil, midtransClient)

	// 1. Non-owner cannot trigger sync (IDOR blocked -> ErrForbidden)
	if _, err := uc.SyncOrderPaymentStatus(context.Background(), orderCode, 999); !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected ErrForbidden for non-owner sync, got %v", err)
	}

	// 2. Owner triggers sync -> queries Midtrans server-to-server -> updates order to PAID
	synced, err := uc.SyncOrderPaymentStatus(context.Background(), orderCode, 50)
	if err != nil {
		t.Fatalf("expected owner sync to succeed, got %v", err)
	}
	if synced.Status != string(domain.OrderStatusPaid) {
		t.Fatalf("expected order status to update to PAID, got %s", synced.Status)
	}
}
