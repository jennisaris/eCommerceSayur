package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	ParentID    *int           `json:"parent_id" gorm:"column:parent_id"`
	Parent      *Category      `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Children    []Category     `json:"children,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Name        string         `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Slug        string         `json:"slug" gorm:"column:slug;type:varchar(120);unique;not null"`
	Icon        string         `json:"icon" gorm:"column:icon;type:varchar(255)"`
	Status      bool           `json:"status" gorm:"column:status;default:true"`
	Description string         `json:"description" gorm:"column:description;type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Category) TableName() string {
	return "categories"
}

type CreateCategoryRequest struct {
	ParentID    *int   `json:"parent_id" validate:"omitempty,gt=0"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"required,min=2,max=120,slug"`
	Icon        string `json:"icon" validate:"omitempty,max=255,safe_url"`
	Status      *bool  `json:"status"`
	Description string `json:"description" validate:"omitempty,max=2000"`
}

type UpdateCategoryRequest struct {
	ParentID    *int   `json:"parent_id" validate:"omitempty,gt=0"`
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Slug        string `json:"slug" validate:"omitempty,min=2,max=120,slug"`
	Icon        string `json:"icon" validate:"omitempty,max=255,safe_url"`
	Status      *bool  `json:"status"`
	Description string `json:"description" validate:"omitempty,max=2000"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id int) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	GetAll(ctx context.Context, onlyActive bool) ([]Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryUseCase interface {
	Create(ctx context.Context, req *CreateCategoryRequest) (*Category, error)
	GetByID(ctx context.Context, id int) (*Category, error)
	GetBySlug(ctx context.Context, slug string) (*Category, error)
	GetAll(ctx context.Context, onlyActive bool) ([]Category, error)
	Update(ctx context.Context, id int, req *UpdateCategoryRequest) (*Category, error)
	Delete(ctx context.Context, id int) error
}
