package usecase

import (
	"context"
	"ecommerce-backend/internal/domain"
	redisrepo "ecommerce-backend/internal/repository/redis"
)

type CategoryUseCase struct {
	categoryRepo domain.CategoryRepository
	catalogCache *redisrepo.CatalogCacheRepository
}

func NewCategoryUseCase(categoryRepo domain.CategoryRepository, catalogCache *redisrepo.CatalogCacheRepository) domain.CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
		catalogCache: catalogCache,
	}
}

func (u *CategoryUseCase) Create(ctx context.Context, req *domain.CreateCategoryRequest) (*domain.Category, error) {
	status := true
	if req.Status != nil {
		status = *req.Status
	}

	category := &domain.Category{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Status:      status,
		Description: req.Description,
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	if u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}

	return category, nil
}

func (u *CategoryUseCase) GetByID(ctx context.Context, id int) (*domain.Category, error) {
	return u.categoryRepo.GetByID(ctx, id)
}

func (u *CategoryUseCase) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return u.categoryRepo.GetBySlug(ctx, slug)
}

func (u *CategoryUseCase) GetAll(ctx context.Context, onlyActive bool) ([]domain.Category, error) {
	return u.categoryRepo.GetAll(ctx, onlyActive)
}

func (u *CategoryUseCase) Update(ctx context.Context, id int, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := u.categoryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.Status != nil {
		category.Status = *req.Status
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := u.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	if u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}

	return category, nil
}

func (u *CategoryUseCase) Delete(ctx context.Context, id int) error {
	err := u.categoryRepo.Delete(ctx, id)
	if err == nil && u.catalogCache != nil {
		_ = u.catalogCache.InvalidateProductCache(ctx)
	}
	return err
}
