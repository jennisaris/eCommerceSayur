package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusProcessing     OrderStatus = "PROCESSING"
	OrderStatusOnDelivery     OrderStatus = "ON_DELIVERY"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// CalculateShippingFee computes dynamic delivery fee based on weight in grams (100% Delivery Model).
// Base fee: Rp 15.000 for first 1,000 grams, Rp 8.000 per additional kg (ceiling).
func CalculateShippingFee(totalWeightGrams int) int64 {
	baseFee := int64(15000)
	if totalWeightGrams <= 1000 {
		return baseFee
	}
	excessGrams := totalWeightGrams - 1000
	additionalKg := (excessGrams + 999) / 1000
	return baseFee + (int64(additionalKg) * 8000)
}

type Order struct {
	ID              int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderCode       string         `json:"order_code" gorm:"column:order_code;type:varchar(60);unique;not null;index"`
	BuyerID         int64          `json:"buyer_id" gorm:"column:buyer_id;not null;index"`
	Buyer           *User          `json:"buyer,omitempty" gorm:"foreignKey:BuyerID;references:ID"`
	Status          string         `json:"status" gorm:"column:status;type:varchar(40);default:'PENDING_PAYMENT';index"`
	RecipientName   string         `json:"recipient_name" gorm:"column:recipient_name;type:varchar(100);not null"`
	RecipientPhone  string         `json:"recipient_phone" gorm:"column:recipient_phone;type:varchar(25);not null"`
	ShippingAddress string         `json:"shipping_address" gorm:"column:shipping_address;type:text;not null"`
	ShippingLat     string         `json:"shipping_lat" gorm:"column:shipping_lat;type:varchar(50)"`
	ShippingLng     string         `json:"shipping_lng" gorm:"column:shipping_lng;type:varchar(50)"`
	ShippingFee     int64          `json:"shipping_fee" gorm:"column:shipping_fee;not null;default:0"`
	TrackingNumber  string         `json:"tracking_number" gorm:"column:tracking_number;type:varchar(100)"`
	TotalAmount     int64          `json:"total_amount" gorm:"column:total_amount;not null"`
	Remarks         string         `json:"remarks" gorm:"column:remarks;type:varchar(255)"`
	Items           []OrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	Payment         *Payment       `json:"payment,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID     int64          `json:"order_id" gorm:"column:order_id;not null;index"`
	ProductID   int64          `json:"product_id" gorm:"column:product_id;not null"`
	Product     *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
	ProductName string         `json:"product_name" gorm:"column:product_name;type:varchar(150);not null"`
	Quantity    int            `json:"quantity" gorm:"column:quantity;not null"`
	Price       int64          `json:"price" gorm:"column:price;not null"`
	Subtotal    int64          `json:"subtotal" gorm:"column:subtotal;not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

// 100% Delivery Checkout DTO: Recipient Name, Phone, Address, Shipping Fee are mandatory
type CheckoutRequest struct {
	RecipientName   string `json:"recipient_name" validate:"required,min=2,max=100"`
	RecipientPhone  string `json:"recipient_phone" validate:"required,phone"`
	ShippingAddress string `json:"shipping_address" validate:"required,min=5,max=500"`
	ShippingLat     string `json:"shipping_lat" validate:"omitempty,max=50,latitude"`
	ShippingLng     string `json:"shipping_lng" validate:"omitempty,max=50,longitude"`
	Remarks         string `json:"remarks" validate:"omitempty,max=255"`
}

type CheckoutEstimateRequest struct {
	ShippingAddress string `json:"shipping_address" validate:"omitempty,max=500"`
	ShippingLat     string `json:"shipping_lat" validate:"omitempty,max=50,latitude"`
	ShippingLng     string `json:"shipping_lng" validate:"omitempty,max=50,longitude"`
}

type CheckoutEstimateResponse struct {
	Subtotal    int64 `json:"subtotal"`
	ShippingFee int64 `json:"shipping_fee"`
	TotalWeight int   `json:"total_weight"` // in grams
	TotalAmount int64 `json:"total_amount"`
}

type UpdateOrderStatusRequest struct {
	Status         string `json:"status" validate:"required,oneof=PROCESSING ON_DELIVERY COMPLETED CANCELLED"`
	TrackingNumber string `json:"tracking_number" validate:"omitempty,max=100"`
}

type OrderFilter struct {
	BuyerID *int64
	Status  string
	Search  string
	Page    int
	Limit   int
}

type OrderRepository interface {
	CreateWithItems(ctx context.Context, order *Order, items []OrderItem) error
	GetByID(ctx context.Context, id int64) (*Order, error)
	GetByCode(ctx context.Context, code string) (*Order, error)
	GetByCodeWithDetails(ctx context.Context, code string) (*Order, error)
	GetAll(ctx context.Context, filter OrderFilter) ([]Order, int64, error)
	UpdateStatus(ctx context.Context, id int64, status string, trackingNumber string) error
	RestockCancelledOrder(ctx context.Context, orderID int64) error
}

type OrderUseCase interface {
	EstimateCheckout(ctx context.Context, userID int64, req *CheckoutEstimateRequest) (*CheckoutEstimateResponse, error)
	CreateOrder(ctx context.Context, userID int64, req *CheckoutRequest) (*Order, error)
	GetOrderByCode(ctx context.Context, userID int64, code string) (*Order, error)
	GetUserOrders(ctx context.Context, userID int64, filter OrderFilter) ([]Order, int64, error)
	AdminGetAllOrders(ctx context.Context, filter OrderFilter) ([]Order, int64, error)
	AdminUpdateOrderStatus(ctx context.Context, orderID int64, req *UpdateOrderStatusRequest) (*Order, error)
}
