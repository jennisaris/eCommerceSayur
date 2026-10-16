package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepository) GetByReceiverID(ctx context.Context, receiverID int64, unreadOnly bool) ([]domain.Notification, error) {
	var notifications []domain.Notification
	query := r.db.WithContext(ctx).Where("receiver_id = ?", receiverID)
	if unreadOnly {
		query = query.Where("read_at IS NULL")
	}

	err := query.Order("created_at DESC").Limit(50).Find(&notifications).Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&domain.Notification{}).
		Where("id = ?", id).
		Update("read_at", &now).Error
}
