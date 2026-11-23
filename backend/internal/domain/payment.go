package domain

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID          int64          `json:"order_id" gorm:"column:order_id;not null;index"`
	Order            *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	UserID           int64          `json:"user_id" gorm:"column:user_id;not null"`
	PaymentMethod    string         `json:"payment_method" gorm:"column:payment_method;type:varchar(50)"`
	PaymentStatus    string         `json:"payment_status" gorm:"column:payment_status;type:varchar(50);default:'pending'"`
	PaymentGatewayID string         `json:"payment_gateway_id" gorm:"column:payment_gateway_id;type:varchar(100);index"`
	GrossAmount      int64          `json:"gross_amount" gorm:"column:gross_amount;not null"`
	PaymentURL       string         `json:"payment_url" gorm:"column:payment_url;type:text"`
	SnapToken        string         `json:"snap_token,omitempty" gorm:"-"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Payment) TableName() string {
	return "payments"
}

type PaymentLog struct {
	ID          int64           `json:"id" gorm:"primaryKey;autoIncrement"`
	PaymentID   int64           `json:"payment_id" gorm:"column:payment_id;not null;index"`
	EventKey    string          `json:"-" gorm:"column:event_key;type:varchar(128);index"`
	Status      string          `json:"status" gorm:"column:status;type:varchar(50);not null"`
	RawResponse json.RawMessage `json:"raw_response" gorm:"column:raw_response;type:jsonb;not null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (PaymentLog) TableName() string {
	return "payment_logs"
}

// MidtransWebhookPayload represents callback from Midtrans
type MidtransWebhookPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"` // This maps to our order_code
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
	SettlementTime    string `json:"settlement_time"`
}

type CreateSnapTransactionResponse struct {
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id int64) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID int64) (*Payment, error)
	GetByGatewayID(ctx context.Context, gatewayID string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	CreateLog(ctx context.Context, log *PaymentLog) error
	GetLogByEventKey(ctx context.Context, eventKey string) (*PaymentLog, error)
}

type PaymentUseCase interface {
	CreateSnapTransaction(ctx context.Context, orderCode string, userID int64) (*CreateSnapTransactionResponse, error)
	HandleMidtransCallback(ctx context.Context, payload []byte) error
	SyncOrderPaymentStatus(ctx context.Context, orderCode string, userID int64) (*Order, error)
}
