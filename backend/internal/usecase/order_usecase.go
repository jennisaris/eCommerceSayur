package usecase

import (
	"context"
	"crypto/rand"
	"ecommerce-backend/internal/domain"
	redisrepo "ecommerce-backend/internal/repository/redis"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type OrderUseCase struct {
	orderRepo   domain.OrderRepository
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
	userRepo    domain.UserRepository
	notifRepo   domain.NotificationRepository
	cartCache   *redisrepo.CartCacheRepository
}

func NewOrderUseCase(
	orderRepo domain.OrderRepository,
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	userRepo domain.UserRepository,
	notifRepo domain.NotificationRepository,
	cartCache *redisrepo.CartCacheRepository,
) domain.OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
		notifRepo:   notifRepo,
		cartCache:   cartCache,
	}
}

// calculateShipping calculates dynamic delivery fee based on weight (100% Delivery Model)
// Base fee: Rp 15.000 for first 1,000 grams, Rp 8.000 per additional kg
func (u *OrderUseCase) calculateShipping(totalWeightGrams int) int64 {
	baseFee := int64(15000)
	if totalWeightGrams <= 1000 {
		return baseFee
	}

	excessGrams := totalWeightGrams - 1000
	additionalKg := (excessGrams + 999) / 1000 // ceiling division
	return baseFee + (int64(additionalKg) * 8000)
}

func (u *OrderUseCase) generateOrderCode() string {
	datePart := time.Now().Format("20060102")
	randomBytes := make([]byte, 3)
	_, _ = rand.Read(randomBytes)
	suffix := strings.ToUpper(hex.EncodeToString(randomBytes))
	return fmt.Sprintf("INV/%s/%s", datePart, suffix)
}

func (u *OrderUseCase) EstimateCheckout(ctx context.Context, userID int64, req *domain.CheckoutEstimateRequest) (*domain.CheckoutEstimateResponse, error) {
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, domain.ErrEmptyCart
	}

	var subtotal int64 = 0
	totalWeight := 0

	for _, item := range cart.Items {
		if item.Product != nil {
			subtotal += item.Product.EffectivePrice() * int64(item.Quantity)
			totalWeight += item.Product.Weight * item.Quantity
		}
	}

	shippingFee := u.calculateShipping(totalWeight)
	totalAmount := subtotal + shippingFee

	return &domain.CheckoutEstimateResponse{
		Subtotal:    subtotal,
		ShippingFee: shippingFee,
		TotalWeight: totalWeight,
		TotalAmount: totalAmount,
	}, nil
}

func (u *OrderUseCase) CreateOrder(ctx context.Context, userID int64, req *domain.CheckoutRequest) (*domain.Order, error) {
	// 1. Mandatory delivery fields check (100% Delivery Model)
	if strings.TrimSpace(req.RecipientName) == "" ||
		strings.TrimSpace(req.RecipientPhone) == "" ||
		strings.TrimSpace(req.ShippingAddress) == "" {
		return nil, domain.ErrInvalidAddress
	}

	// 2. Fetch current user cart
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(cart.Items) == 0 {
		return nil, domain.ErrEmptyCart
	}

	// 3. Prepare order items & calculate total weight for shipping fee
	var orderItems []domain.OrderItem
	totalWeight := 0

	for _, item := range cart.Items {
		if item.Quantity <= 0 {
			continue
		}
		if item.Product != nil {
			totalWeight += item.Product.Weight * item.Quantity
		}
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	if len(orderItems) == 0 {
		return nil, domain.ErrEmptyCart
	}

	shippingFee := u.calculateShipping(totalWeight)

	order := &domain.Order{
		OrderCode:       u.generateOrderCode(),
		BuyerID:         userID,
		Status:          string(domain.OrderStatusPendingPayment),
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		ShippingAddress: req.ShippingAddress,
		ShippingLat:     req.ShippingLat,
		ShippingLng:     req.ShippingLng,
		ShippingFee:     shippingFee,
		Remarks:         req.Remarks,
	}

	// 4. Execute atomic transaction (lock stock, price snapshot, create order & payment)
	if err := u.orderRepo.CreateWithItems(ctx, order, orderItems); err != nil {
		return nil, err
	}

	// 5. Invalidate cart cache
	if u.cartCache != nil {
		_ = u.cartCache.InvalidateCart(ctx, userID)
	}

	// 6. Send customer notification
	if u.notifRepo != nil {
		_ = u.notifRepo.Create(ctx, &domain.Notification{
			ReceiverID:       userID,
			NotificationType: "ORDER_STATUS",
			Subject:          fmt.Sprintf("Pesanan Dibuat (%s)", order.OrderCode),
			Message:          fmt.Sprintf("Pesanan Anda dengan kode %s telah dibuat. Silakan selesaikan pembayaran sebelum batas waktu berakhir.", order.OrderCode),
		})
	}

	return order, nil
}

func (u *OrderUseCase) GetOrderByCode(ctx context.Context, userID int64, code string, isAdmin bool) (*domain.Order, error) {
	order, err := u.orderRepo.GetByCodeWithDetails(ctx, code)
	if err != nil {
		return nil, err
	}

	if !isAdmin && order.BuyerID != userID {
		return nil, domain.ErrForbidden
	}

	return order, nil
}

func (u *OrderUseCase) GetUserOrders(ctx context.Context, userID int64, filter domain.OrderFilter) ([]domain.Order, int64, error) {
	filter.BuyerID = &userID
	return u.orderRepo.GetAll(ctx, filter)
}

func (u *OrderUseCase) AdminGetAllOrders(ctx context.Context, filter domain.OrderFilter) ([]domain.Order, int64, error) {
	return u.orderRepo.GetAll(ctx, filter)
}

func (u *OrderUseCase) AdminUpdateOrderStatus(ctx context.Context, orderID int64, req *domain.UpdateOrderStatusRequest) (*domain.Order, error) {
	order, err := u.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// If transitioning to CANCELLED, restock automatically
	if req.Status == string(domain.OrderStatusCancelled) && order.Status != string(domain.OrderStatusCancelled) {
		if err := u.orderRepo.RestockCancelledOrder(ctx, orderID); err != nil {
			return nil, err
		}
	} else {
		if err := u.orderRepo.UpdateStatus(ctx, orderID, req.Status, req.TrackingNumber); err != nil {
			return nil, err
		}
	}

	// Send customer notification
	if u.notifRepo != nil {
		subject := fmt.Sprintf("Status Pesanan Diperbarui: %s", req.Status)
		msg := fmt.Sprintf("Pesanan %s kini berstatus %s.", order.OrderCode, req.Status)
		if req.TrackingNumber != "" {
			msg += fmt.Sprintf(" Nomor Resi Pengiriman: %s", req.TrackingNumber)
		}
		_ = u.notifRepo.Create(ctx, &domain.Notification{
			ReceiverID:       order.BuyerID,
			NotificationType: "ORDER_STATUS",
			Subject:          subject,
			Message:          msg,
		})
	}

	return u.orderRepo.GetByID(ctx, orderID)
}
