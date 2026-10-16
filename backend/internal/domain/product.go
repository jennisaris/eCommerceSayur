package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID   int            `json:"category_id" gorm:"column:category_id;not null"`
	Category     *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	ParentID     *int64         `json:"parent_id" gorm:"column:parent_id"`
	Variants     []Product      `json:"variants,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Name         string         `json:"name" gorm:"column:name;type:varchar(150);not null"`
	Slug         string         `json:"slug" gorm:"column:slug;type:varchar(180);unique;not null"`
	Image        string         `json:"image" gorm:"column:image;type:text"`
	Description  string         `json:"description" gorm:"column:description;type:text"`
	RegularPrice int64          `json:"regular_price" gorm:"column:regular_price;not null"`
	SalePrice    int64          `json:"sale_price" gorm:"column:sale_price;default:0"`
	Unit         string         `json:"unit" gorm:"column:unit;type:varchar(30);default:'pcs'"`
	Weight       int            `json:"weight" gorm:"column:weight;not null;default:1000"` // in grams
	Stock        int            `json:"stock" gorm:"column:stock;not null;default:0"`
	VariantName  string         `json:"variant_name" gorm:"column:variant_name;type:varchar(50)"`
	Status       string         `json:"status" gorm:"column:status;type:varchar(30);default:'ACTIVE'"` // ACTIVE, INACTIVE, OUT_OF_STOCK
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
}

func (Product) TableName() string {
	return "products"
}

// EffectivePrice returns SalePrice if > 0, else RegularPrice
func (p *Product) EffectivePrice() int64 {
	if p.SalePrice > 0 && p.SalePrice < p.RegularPrice {
		return p.SalePrice
	}
	return p.RegularPrice
}

type ProductFilter struct {
	CategoryID *int
	Search     string
	Status     string
	MinPrice   *int64
	MaxPrice   *int64
	Page       int
	Limit      int
	Sort       string // "price_asc", "price_desc", "latest"
}

type CreateProductRequest struct {
	CategoryID   int    `json:"category_id" validate:"required"`
	ParentID     *int64 `json:"parent_id"`
	Name         string `json:"name" validate:"required,min=2,max=150"`
	Slug         string `json:"slug" validate:"required,min=2,max=180"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	RegularPrice int64  `json:"regular_price" validate:"required,gt=0"`
	SalePrice    int64  `json:"sale_price" validate:"gte=0"`
	Unit         string `json:"unit"`
	Weight       int    `json:"weight" validate:"required,gt=0"`
	Stock        int    `json:"stock" validate:"gte=0"`
	VariantName  string `json:"variant_name"`
	Status       string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE OUT_OF_STOCK"`
}

type UpdateProductRequest struct {
	CategoryID   *int   `json:"category_id"`
	ParentID     *int64 `json:"parent_id"`
	Name         string `json:"name" validate:"omitempty,min=2,max=150"`
	Slug         string `json:"slug" validate:"omitempty,min=2,max=180"`
	Image        string `json:"image"`
	Description  string `json:"description"`
	RegularPrice *int64 `json:"regular_price" validate:"omitempty,gt=0"`
	SalePrice    *int64 `json:"sale_price" validate:"omitempty,gte=0"`
	Unit         string `json:"unit"`
	Weight       *int   `json:"weight" validate:"omitempty,gt=0"`
	Stock        *int   `json:"stock" validate:"omitempty,gte=0"`
	VariantName  string `json:"variant_name"`
	Status       string `json:"status" validate:"omitempty,oneof=ACTIVE INACTIVE OUT_OF_STOCK"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	GetAll(ctx context.Context, filter ProductFilter) ([]Product, int64, error)
	Update(ctx context.Context, product *Product) error
	UpdateStock(ctx context.Context, tx *gorm.DB, productID int64, delta int) error
	Delete(ctx context.Context, id int64) error
}

type ProductUseCase interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetBySlug(ctx context.Context, slug string) (*Product, error)
	GetAll(ctx context.Context, filter ProductFilter) ([]Product, int64, error)
	Update(ctx context.Context, id int64, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id int64) error
}
