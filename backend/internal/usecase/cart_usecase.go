package usecase

import (
	"context"
	"ecommerce-backend/internal/domain"
	redisrepo "ecommerce-backend/internal/repository/redis"
	"errors"
	"fmt"
	"time"
)

type CartUseCase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
	cartCache   *redisrepo.CartCacheRepository
}

func NewCartUseCase(
	cartRepo domain.CartRepository,
	productRepo domain.ProductRepository,
	cartCache *redisrepo.CartCacheRepository,
) domain.CartUseCase {
	return &CartUseCase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		cartCache:   cartCache,
	}
}

func (u *CartUseCase) calculateSummary(cart *domain.Cart) *domain.CartSummaryResponse {
	totalItems := 0
	totalWeight := 0
	var totalPrice int64 = 0

	for _, item := range cart.Items {
		totalItems += item.Quantity
		if item.Product != nil {
			totalWeight += item.Product.Weight * item.Quantity
			totalPrice += item.Product.EffectivePrice() * int64(item.Quantity)
		}
	}

	return &domain.CartSummaryResponse{
		Cart:        cart,
		TotalItems:  totalItems,
		TotalWeight: totalWeight,
		TotalPrice:  totalPrice,
	}
}

func (u *CartUseCase) GetCart(ctx context.Context, userID int64) (*domain.CartSummaryResponse, error) {
	// Try Redis cache
	if u.cartCache != nil {
		if cached, err := u.cartCache.GetCart(ctx, userID); err == nil && cached != nil {
			return cached, nil
		}
	}

	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := u.calculateSummary(cart)

	// Cache for 30 minutes
	if u.cartCache != nil {
		_ = u.cartCache.SetCart(ctx, userID, summary, 30*time.Minute)
	}

	return summary, nil
}

func (u *CartUseCase) AddToCart(ctx context.Context, userID int64, req *domain.AddToCartRequest) (*domain.CartSummaryResponse, error) {
	product, err := u.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", domain.ErrNotFound)
	}

	if product.Status != "ACTIVE" {
		return nil, domain.ErrProductInactive
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("requested quantity (%d) exceeds available stock (%d): %w",
			req.Quantity, product.Stock, domain.ErrInsufficientStock)
	}

	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	item := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := u.cartRepo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	if u.cartCache != nil {
		_ = u.cartCache.InvalidateCart(ctx, userID)
	}

	return u.GetCart(ctx, userID)
}

func (u *CartUseCase) UpdateItem(ctx context.Context, userID int64, productID int64, req *domain.UpdateCartItemRequest) (*domain.CartSummaryResponse, error) {
	if req.Quantity <= 0 {
		return u.RemoveItem(ctx, userID, productID)
	}

	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	product, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("requested quantity (%d) exceeds available stock (%d): %w",
			req.Quantity, product.Stock, domain.ErrInsufficientStock)
	}

	if err := u.cartRepo.UpdateItemQuantity(ctx, cart.ID, productID, req.Quantity); err != nil {
		return nil, err
	}

	if u.cartCache != nil {
		_ = u.cartCache.InvalidateCart(ctx, userID)
	}

	return u.GetCart(ctx, userID)
}

func (u *CartUseCase) RemoveItem(ctx context.Context, userID int64, productID int64) (*domain.CartSummaryResponse, error) {
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := u.cartRepo.RemoveItem(ctx, cart.ID, productID); err != nil {
		return nil, err
	}

	if u.cartCache != nil {
		_ = u.cartCache.InvalidateCart(ctx, userID)
	}

	return u.GetCart(ctx, userID)
}

func (u *CartUseCase) ClearCart(ctx context.Context, userID int64) error {
	cart, err := u.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}

	if err := u.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return err
	}

	if u.cartCache != nil {
		_ = u.cartCache.InvalidateCart(ctx, userID)
	}

	return nil
}
