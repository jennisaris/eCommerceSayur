package domain

import (
	"context"
	"time"
)

type Cart struct {
	ID        int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64      `json:"user_id" gorm:"column:user_id;unique;not null"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID;references:ID"`
	CreatedAt time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    int64     `json:"cart_id" gorm:"column:cart_id;not null"`
	ProductID int64     `json:"product_id" gorm:"column:product_id;not null"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int       `json:"quantity" gorm:"column:quantity;not null;default:1"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

type AddToCartRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"required,min=1,max=1000"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"gte=0,lte=1000"`
}

type CartSummaryResponse struct {
	Cart        *Cart `json:"cart"`
	TotalItems  int   `json:"total_items"`
	TotalWeight int   `json:"total_weight"` // grams
	TotalPrice  int64 `json:"total_price"`  // IDR
}

type CartRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*Cart, error)
	Create(ctx context.Context, cart *Cart) error
	AddItem(ctx context.Context, item *CartItem) error
	UpdateItemQuantity(ctx context.Context, cartID int64, productID int64, quantity int) error
	RemoveItem(ctx context.Context, cartID int64, productID int64) error
	ClearCart(ctx context.Context, cartID int64) error
}

type CartUseCase interface {
	GetCart(ctx context.Context, userID int64) (*CartSummaryResponse, error)
	AddToCart(ctx context.Context, userID int64, req *AddToCartRequest) (*CartSummaryResponse, error)
	UpdateItem(ctx context.Context, userID int64, productID int64, req *UpdateCartItemRequest) (*CartSummaryResponse, error)
	RemoveItem(ctx context.Context, userID int64, productID int64) (*CartSummaryResponse, error)
	ClearCart(ctx context.Context, userID int64) error
}
