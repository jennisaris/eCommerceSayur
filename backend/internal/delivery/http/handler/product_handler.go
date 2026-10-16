package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUseCase domain.ProductUseCase
	validate       *validator.Validate
}

func NewProductHandler(productUseCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		validate:       validator.New(),
	}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")
	status := c.Query("status", "ACTIVE") // Default to ACTIVE for public store
	sort := c.Query("sort", "latest")

	var categoryID *int
	if catStr := c.Query("category_id"); catStr != "" {
		if id, err := strconv.Atoi(catStr); err == nil && id > 0 {
			categoryID = &id
		}
	}

	var minPrice *int64
	if minStr := c.Query("min_price"); minStr != "" {
		if val, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			minPrice = &val
		}
	}

	var maxPrice *int64
	if maxStr := c.Query("max_price"); maxStr != "" {
		if val, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			maxPrice = &val
		}
	}

	// Admin can pass status="all" to see all products including inactive
	if c.Query("status") == "all" {
		status = ""
	}

	filter := domain.ProductFilter{
		CategoryID: categoryID,
		Search:     search,
		Status:     status,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       page,
		Limit:      limit,
		Sort:       sort,
	}

	products, total, err := h.productUseCase.GetAll(c.Context(), filter)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Paginated(c, "Products retrieved successfully", products, page, limit, total)
}

func (h *ProductHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	product, err := h.productUseCase.GetBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, "Product details retrieved successfully", product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := h.productUseCase.Create(c.Context(), &req)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Created(c, "Product created successfully", product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req domain.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	product, err := h.productUseCase.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	if err := h.productUseCase.Delete(c.Context(), id); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Product deleted successfully", nil)
}
