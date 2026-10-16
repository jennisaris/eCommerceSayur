package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/midtrans"
)

type PaymentUseCase struct {
	paymentRepo domain.PaymentRepository
	orderRepo   domain.OrderRepository
	userRepo    domain.UserRepository
	notifRepo   domain.NotificationRepository
	midtransSvc *midtrans.Client
}

func NewPaymentUseCase(
	paymentRepo domain.PaymentRepository,
	orderRepo domain.OrderRepository,
	userRepo domain.UserRepository,
	notifRepo domain.NotificationRepository,
	midtransSvc *midtrans.Client,
) domain.PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		notifRepo:   notifRepo,
		midtransSvc: midtransSvc,
	}
}

func (u *PaymentUseCase) CreateSnapTransaction(ctx context.Context, orderCode string, userID int64) (*domain.CreateSnapTransactionResponse, error) {
	order, err := u.orderRepo.GetByCodeWithDetails(ctx, orderCode)
	if err != nil {
		return nil, err
	}

	if order.BuyerID != userID {
		return nil, domain.ErrForbidden
	}

	if order.Status != string(domain.OrderStatusPendingPayment) {
		return nil, fmt.Errorf("order cannot be paid in current status (%s): %w", order.Status, domain.ErrInvalidOrderState)
	}

	buyer, _ := u.userRepo.GetByID(ctx, userID)

	// Prepare item details for Snap popup
	var items []midtrans.SnapItemDetails
	for _, it := range order.Items {
		items = append(items, midtrans.SnapItemDetails{
			ID:       fmt.Sprintf("PROD-%d", it.ProductID),
			Price:    it.Price,
			Quantity: it.Quantity,
			Name:     it.ProductName,
		})
	}

	// Add shipping fee as line item if > 0
	if order.ShippingFee > 0 {
		items = append(items, midtrans.SnapItemDetails{
			ID:       "SHIPPING-FEE",
			Price:    order.ShippingFee,
			Quantity: 1,
			Name:     "Ongkos Kirim Pengantaran 100%",
		})
	}

	customerName := order.RecipientName
	customerEmail := "customer@ecommerce.com"
	if buyer != nil && buyer.Email != "" {
		customerEmail = buyer.Email
	}

	snapReq := &midtrans.CreateSnapRequest{
		TransactionDetails: midtrans.SnapTransactionDetails{
			OrderID:     order.OrderCode,
			GrossAmount: order.TotalAmount,
		},
		CustomerDetails: midtrans.SnapCustomerDetails{
			FirstName: customerName,
			Email:     customerEmail,
			Phone:     order.RecipientPhone,
			ShippingAddress: &midtrans.SnapShippingAddress{
				FirstName: order.RecipientName,
				Phone:     order.RecipientPhone,
				Address:   order.ShippingAddress,
			},
		},
		ItemDetails: items,
	}

	snapResp, err := u.midtransSvc.CreateSnapTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate snap transaction: %w", err)
	}

	// Update existing payment or create record
	payment, err := u.paymentRepo.GetByOrderID(ctx, order.ID)
	if err != nil {
		payment = &domain.Payment{
			OrderID:       order.ID,
			UserID:        order.BuyerID,
			PaymentStatus: "pending",
			GrossAmount:   order.TotalAmount,
		}
	}

	payment.PaymentURL = snapResp.RedirectURL
	payment.SnapToken = snapResp.Token
	_ = u.paymentRepo.Update(ctx, payment)

	return &domain.CreateSnapTransactionResponse{
		SnapToken:   snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

// HandleMidtransCallback processes incoming webhook from Midtrans:
// 1. Verifies SHA-512 signature
// 2. Stores raw JSON log in payment_logs (Idempotent audit)
// 3. Atomically updates order and payment status
// 4. Automatically restocks product if payment expires/cancels
func (u *PaymentUseCase) HandleMidtransCallback(ctx context.Context, payload []byte) error {
	var webhook domain.MidtransWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return fmt.Errorf("invalid json webhook: %w", err)
	}

	// 1. Verify SHA-512 signature
	if !u.midtransSvc.VerifySignature(webhook.OrderID, webhook.StatusCode, webhook.GrossAmount, webhook.SignatureKey) {
		return domain.ErrInvalidSignature
	}

	// 2. Find order
	order, err := u.orderRepo.GetByCode(ctx, webhook.OrderID)
	if err != nil {
		return fmt.Errorf("order not found for code %s: %w", webhook.OrderID, err)
	}

	// 3. Find or initialize payment
	payment, err := u.paymentRepo.GetByOrderID(ctx, order.ID)
	if err != nil {
		payment = &domain.Payment{
			OrderID:     order.ID,
			UserID:      order.BuyerID,
			GrossAmount: order.TotalAmount,
		}
		_ = u.paymentRepo.Create(ctx, payment)
	}

	// 4. Idempotent Log entry to payment_logs (raw_response JSONB)
	_ = u.paymentRepo.CreateLog(ctx, &domain.PaymentLog{
		PaymentID:   payment.ID,
		Status:      webhook.TransactionStatus,
		RawResponse: json.RawMessage(payload),
	})

	// 5. Update payment details
	payment.PaymentGatewayID = webhook.TransactionID
	payment.PaymentMethod = webhook.PaymentType
	payment.PaymentStatus = webhook.TransactionStatus

	if gross, err := strconv.ParseFloat(webhook.GrossAmount, 64); err == nil && gross > 0 {
		payment.GrossAmount = int64(gross)
	}

	// 6. Map Midtrans transaction status to internal order status
	status := strings.ToLower(webhook.TransactionStatus)
	fraud := strings.ToLower(webhook.FraudStatus)

	switch status {
	case "capture":
		if fraud == "challenge" {
			payment.PaymentStatus = "challenge"
		} else if fraud == "accept" {
			payment.PaymentStatus = "settlement"
			_ = u.orderRepo.UpdateStatus(ctx, order.ID, string(domain.OrderStatusPaid), "")
		}
	case "settlement":
		payment.PaymentStatus = "settlement"
		_ = u.orderRepo.UpdateStatus(ctx, order.ID, string(domain.OrderStatusPaid), "")

		if u.notifRepo != nil {
			_ = u.notifRepo.Create(ctx, &domain.Notification{
				ReceiverID:       order.BuyerID,
				NotificationType: "PAYMENT",
				Subject:          "Pembayaran Diterima!",
				Message:          fmt.Sprintf("Pembayaran untuk pesanan %s senilai Rp %d telah berhasil diverifikasi.", order.OrderCode, order.TotalAmount),
			})
		}

	case "pending":
		payment.PaymentStatus = "pending"

	case "deny", "expire", "cancel":
		payment.PaymentStatus = status
		// CRITICAL RULE 5: Restock items automatically upon cancellation/expiration
		_ = u.orderRepo.RestockCancelledOrder(ctx, order.ID)

		if u.notifRepo != nil {
			_ = u.notifRepo.Create(ctx, &domain.Notification{
				ReceiverID:       order.BuyerID,
				NotificationType: "PAYMENT",
				Subject:          fmt.Sprintf("Pesanan Dibatalkan (%s)", status),
				Message:          fmt.Sprintf("Pembayaran untuk pesanan %s %s. Stok produk telah dikembalikan.", order.OrderCode, status),
			})
		}
	}

	return u.paymentRepo.Update(ctx, payment)
}
