package postgres

import (
	"context"
	"ecommerce-backend/internal/domain"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDForUpdate(ctx context.Context, tx *gorm.DB, id int64) (*domain.Product, error) {
	var product domain.Product
	db := r.db
	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Variants").
		Where("slug = ?", slug).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{}).Preload("Category")

	if filter.CategoryID != nil && *filter.CategoryID > 0 {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Search != "" {
		searchTerm := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", searchTerm, searchTerm)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.MinPrice != nil {
		query = query.Where("(CASE WHEN sale_price > 0 THEN sale_price ELSE regular_price END) >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("(CASE WHEN sale_price > 0 THEN sale_price ELSE regular_price END) <= ?", *filter.MaxPrice)
	}

	// Count total matching items before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	switch filter.Sort {
	case "price_asc":
		query = query.Order("(CASE WHEN sale_price > 0 THEN sale_price ELSE regular_price END) ASC")
	case "price_desc":
		query = query.Order("(CASE WHEN sale_price > 0 THEN sale_price ELSE regular_price END) DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductRepository) UpdateStock(ctx context.Context, tx *gorm.DB, productID int64, delta int) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	return db.WithContext(ctx).
		Model(&domain.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", delta)).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}
