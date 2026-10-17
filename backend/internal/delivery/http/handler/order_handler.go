package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderUseCase domain.OrderUseCase
	validate     *validator.Validate
}

func NewOrderHandler(orderUseCase domain.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		validate:     validator.New(),
	}
}

func (h *OrderHandler) EstimateCheckout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req domain.CheckoutEstimateRequest
	_ = c.BodyParser(&req)

	estimate, err := h.orderUseCase.EstimateCheckout(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyCart) {
			return response.Error(c, fiber.StatusBadRequest, "Keranjang belanja masih kosong")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Estimasi biaya pengiriman berhasil", estimate)
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req domain.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	order, err := h.orderUseCase.CreateOrder(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyCart) {
			return response.Error(c, fiber.StatusBadRequest, "Keranjang belanja masih kosong")
		}
		if errors.Is(err, domain.ErrInsufficientStock) || errors.Is(err, domain.ErrProductInactive) {
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidAddress) {
			return response.Error(c, fiber.StatusBadRequest, "Nama penerima, no telepon, dan alamat pengiriman wajib diisi lengkap")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Created(c, "Pesanan berhasil dibuat", order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")

	filter := domain.OrderFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := h.orderUseCase.GetUserOrders(c.Context(), userID, filter)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Paginated(c, "Orders retrieved successfully", orders, page, limit, total)
}

func (h *OrderHandler) GetOrderByCode(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	code := c.Params("code")
	if code == "" {
		code = strings.TrimPrefix(c.Params("*"), "/")
	}

	// Check if user is Super Admin
	isAdmin := false
	if roles, ok := c.Locals("roles").([]string); ok {
		for _, r := range roles {
			if r == "Super Admin" {
				isAdmin = true
				break
			}
		}
	}

	order, err := h.orderUseCase.GetOrderByCode(c.Context(), userID, code, isAdmin)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Pesanan tidak ditemukan")
		}
		if errors.Is(err, domain.ErrForbidden) {
			return response.Error(c, fiber.StatusForbidden, "Akses pesanan ditolak")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Detail pesanan berhasil diambil", order)
}
