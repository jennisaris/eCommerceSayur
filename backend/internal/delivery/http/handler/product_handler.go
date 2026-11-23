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

type ProductHandler struct {
	productUseCase domain.ProductUseCase
	validate       *validator.Validate
}

func NewProductHandler(productUseCase domain.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		validate:       pkgvalidate.New(),
	}
}

func parseProductFilter(c *fiber.Ctx, allowInactive bool) (domain.ProductFilter, int, int, error) {
	page, limit := parsePagination(c)
	search := strings.TrimSpace(c.Query("search", ""))
	if len(search) > 100 {
		search = search[:100]
	}

	status := "ACTIVE"
	if allowInactive {
		rawStatus := strings.TrimSpace(c.Query("status", ""))
		switch rawStatus {
		case "", "all":
			status = ""
		case "ACTIVE", "INACTIVE", "OUT_OF_STOCK":
			status = rawStatus
		default:
			return domain.ProductFilter{}, 0, 0, errors.New("Invalid product status filter")
		}
	}

	sort := strings.TrimSpace(c.Query("sort", "latest"))
	switch sort {
	case "", "latest", "price_asc", "price_desc":
		if sort == "" {
			sort = "latest"
		}
	default:
		return domain.ProductFilter{}, 0, 0, errors.New("Invalid sort parameter")
	}

	var categoryID *int
	if catStr := strings.TrimSpace(c.Query("category_id")); catStr != "" {
		id, err := strconv.Atoi(catStr)
		if err != nil || id <= 0 {
			return domain.ProductFilter{}, 0, 0, errors.New("Invalid category_id parameter")
		}
		categoryID = &id
	}

	var minPrice *int64
	if minStr := strings.TrimSpace(c.Query("min_price")); minStr != "" {
		val, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil || val < 0 {
			return domain.ProductFilter{}, 0, 0, errors.New("Invalid min_price parameter")
		}
		minPrice = &val
	}

	var maxPrice *int64
	if maxStr := strings.TrimSpace(c.Query("max_price")); maxStr != "" {
		val, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil || val < 0 {
			return domain.ProductFilter{}, 0, 0, errors.New("Invalid max_price parameter")
		}
		maxPrice = &val
	}
	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		return domain.ProductFilter{}, 0, 0, errors.New("min_price cannot exceed max_price")
	}

	return domain.ProductFilter{
		CategoryID: categoryID,
		Search:     search,
		Status:     status,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Page:       page,
		Limit:      limit,
		Sort:       sort,
	}, page, limit, nil
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	filter, page, limit, err := parseProductFilter(c, false)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	products, total, err := h.productUseCase.GetAll(c.Context(), filter)
	if err != nil {
		return response.Internal(c, err)
	}

	return response.Paginated(c, "Products retrieved successfully", products, page, limit, total)
}

func (h *ProductHandler) AdminGetAll(c *fiber.Ctx) error {
	filter, page, limit, err := parseProductFilter(c, true)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	products, total, err := h.productUseCase.GetAll(c.Context(), filter)
	if err != nil {
		return response.Internal(c, err)
	}

	return response.Paginated(c, "Products retrieved successfully", products, page, limit, total)
}

func (h *ProductHandler) GetBySlug(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" || len(slug) > 180 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product slug")
	}
	product, err := h.productUseCase.GetBySlug(c.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found")
		}
		return response.Internal(c, err)
	}
	if product.Status == "INACTIVE" {
		return response.Error(c, fiber.StatusNotFound, "Product not found")
	}
	return response.Success(c, "Product details retrieved successfully", product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req domain.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	product, err := h.productUseCase.Create(c.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) || errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return response.Internal(c, err)
	}

	return response.Created(c, "Product created successfully", product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req domain.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	product, err := h.productUseCase.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Product not found")
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, err.Error())
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	if err := h.productUseCase.Delete(c.Context(), id); err != nil {
		return response.Internal(c, err)
	}

	return response.Success(c, "Product deleted successfully", nil)
}
