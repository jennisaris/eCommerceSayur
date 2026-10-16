package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ReceiverID       int64          `json:"receiver_id" gorm:"column:receiver_id;not null;index"`
	NotificationType string         `json:"notification_type" gorm:"column:notification_type;type:varchar(50);not null"` // ORDER_STATUS, PAYMENT, PROMO, SYSTEM
	Subject          string         `json:"subject" gorm:"column:subject;type:varchar(150);not null"`
	Message          string         `json:"message" gorm:"column:message;type:text;not null"`
	SentAt           time.Time      `json:"sent_at" gorm:"column:sent_at;autoCreateTime"`
	ReadAt           *time.Time     `json:"read_at" gorm:"column:read_at"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Notification) TableName() string {
	return "notifications"
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByReceiverID(ctx context.Context, receiverID int64, unreadOnly bool) ([]Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
}
