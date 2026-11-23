package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	pkgvalidate "ecommerce-backend/internal/pkg/validate"
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryUseCase domain.CategoryUseCase
	validate        *validator.Validate
}

func NewCategoryHandler(categoryUseCase domain.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
		validate:        pkgvalidate.New(),
	}
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	categories, err := h.categoryUseCase.GetAll(c.Context(), true)
	if err != nil {
		return response.Internal(c, err)
	}
	return response.Success(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) AdminGetAll(c *fiber.Ctx) error {
	onlyActive := c.Query("active", "false") == "true"
	categories, err := h.categoryUseCase.GetAll(c.Context(), onlyActive)
	if err != nil {
		return response.Internal(c, err)
	}
	return response.Success(c, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) GetBySlug(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" || len(slug) > 120 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid category slug")
	}
	category, err := h.categoryUseCase.GetBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Category not found")
		}
		return response.Internal(c, err)
	}
	if !category.Status {
		return response.Error(c, fiber.StatusNotFound, "Category not found")
	}
	return response.Success(c, "Category retrieved successfully", category)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	category, err := h.categoryUseCase.Create(c.Context(), &req)
	if err != nil {
		return response.Internal(c, err)
	}

	return response.Created(c, "Category created successfully", category)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	var req domain.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	category, err := h.categoryUseCase.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Category not found")
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Category updated successfully", category)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	if err := h.categoryUseCase.Delete(c.Context(), id); err != nil {
		return response.Internal(c, err)
	}

	return response.Success(c, "Category deleted successfully", nil)
}
