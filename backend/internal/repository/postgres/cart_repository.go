package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) domain.CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetByUserID(ctx context.Context, userID int64) (*domain.Cart, error) {
	var cart domain.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items.Product.Category").
		First(&cart).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create empty cart for this user
			newCart := domain.Cart{UserID: userID}
			if err := r.db.WithContext(ctx).Create(&newCart).Error; err != nil {
				return nil, err
			}
			newCart.Items = []domain.CartItem{}
			return &newCart, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *CartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *CartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	var existing domain.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).
		First(&existing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.WithContext(ctx).Create(item).Error
		}
		return err
	}

	// If already in cart, increment quantity
	existing.Quantity += item.Quantity
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *CartRepository) UpdateItemQuantity(ctx context.Context, cartID int64, productID int64, quantity int) error {
	if quantity <= 0 {
		return r.RemoveItem(ctx, cartID, productID)
	}

	return r.db.WithContext(ctx).
		Model(&domain.CartItem{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(ctx context.Context, cartID int64, productID int64) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&domain.CartItem{}).Error
}

func (r *CartRepository) ClearCart(ctx context.Context, cartID int64) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&domain.CartItem{}).Error
}
