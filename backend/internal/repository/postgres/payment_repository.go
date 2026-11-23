package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Preload("Order").
		First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID int64) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetByGatewayID(ctx context.Context, gatewayID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).
		Where("payment_gateway_id = ?", gatewayID).
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *PaymentRepository) CreateLog(ctx context.Context, log *domain.PaymentLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *PaymentRepository) GetLogByEventKey(ctx context.Context, eventKey string) (*domain.PaymentLog, error) {
	var paymentLog domain.PaymentLog
	err := r.db.WithContext(ctx).
		Where("event_key = ?", eventKey).
		First(&paymentLog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &paymentLog, nil
}
