package usecase

import (
	"context"
	"ecommerce-backend/internal/domain"
	redisrepo "ecommerce-backend/internal/repository/redis"
	"fmt"
	"time"
)

type ProductUseCase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
	catalogCache *redisrepo.CatalogCacheRepository
}

func NewProductUseCase(
	productRepo domain.ProductRepository,
	categoryRepo domain.CategoryRepository,
	catalogCache *redisrepo.CatalogCacheRepository,
) domain.ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		catalogCache: catalogCache,
	}
}

func (u *ProductUseCase) Create(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	if req.SalePrice > 0 && req.SalePrice > req.RegularPrice {
		return nil, fmt.Errorf("harga diskon tidak boleh melebihi harga normal: %w", domain.ErrInvalidInput)
	}

	// Validate category exists
	if _, err := u.categoryRepo.GetByID(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	status := "ACTIVE"
	if req.Status != "" {
		status = req.Status
	}

	unit := "pcs"
	if req.Unit != "" {
		unit = req.Unit
	}

	product := &domain.Product{
		CategoryID:   req.CategoryID,
		ParentID:     req.ParentID,
		Name:         req.Name,
		Slug:         req.Slug,
		Image:        req.Image,
		Description:  req.Description,
		RegularPrice: req.RegularPrice,
		SalePrice:    req.SalePrice,
		Unit:         unit,
		Weight:       req.Weight,
		Stock:        req.Stock,
		VariantName:  req.VariantName,
		Status:       status,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	if u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}

	return product, nil
}

func (u *ProductUseCase) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *ProductUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	// Try Redis cache
	if u.catalogCache != nil {
		if cached, err := u.catalogCache.GetProductBySlug(ctx, slug); err == nil && cached != nil {
			return cached, nil
		}
	}

	product, err := u.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	// Cache in Redis for 15 minutes
	if u.catalogCache != nil {
		_ = u.catalogCache.SetProductBySlug(ctx, slug, product, 15*time.Minute)
	}

	return product, nil
}

func (u *ProductUseCase) GetAll(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int64, error) {
	var cacheKey string
	if u.catalogCache != nil {
		cacheKey = u.catalogCache.GenerateProductFilterKey(filter)
		if cached, err := u.catalogCache.GetProducts(ctx, cacheKey); err == nil && cached != nil {
			return cached.Products, cached.Total, nil
		}
	}

	products, total, err := u.productRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Cache in Redis for 5 minutes
	if u.catalogCache != nil && cacheKey != "" {
		_ = u.catalogCache.SetProducts(ctx, cacheKey, products, total, 5*time.Minute)
	}

	return products, total, nil
}

func (u *ProductUseCase) Update(ctx context.Context, id int64, req *domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := u.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != nil {
		if _, err := u.categoryRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, err
		}
		product.CategoryID = *req.CategoryID
	}

	if req.ParentID != nil {
		product.ParentID = req.ParentID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	if req.Image != "" {
		product.Image = req.Image
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.RegularPrice != nil {
		product.RegularPrice = *req.RegularPrice
	}
	if req.SalePrice != nil {
		product.SalePrice = *req.SalePrice
	}
	if product.SalePrice > 0 && product.SalePrice > product.RegularPrice {
		return nil, fmt.Errorf("harga diskon tidak boleh melebihi harga normal: %w", domain.ErrInvalidInput)
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.Weight != nil {
		product.Weight = *req.Weight
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
		if product.Stock == 0 {
			product.Status = "OUT_OF_STOCK"
		}
	}
	if req.VariantName != "" {
		product.VariantName = req.VariantName
	}
	if req.Status != "" {
		product.Status = req.Status
	}

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	if u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}

	return product, nil
}

func (u *ProductUseCase) Delete(ctx context.Context, id int64) error {
	err := u.productRepo.Delete(ctx, id)
	if err == nil && u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}
	return err
}
