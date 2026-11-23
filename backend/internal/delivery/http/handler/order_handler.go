package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"
	pkgvalidate "ecommerce-backend/internal/pkg/validate"
	"errors"
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
		validate:     pkgvalidate.New(),
	}
}

func isValidOrderStatusFilter(status string) bool {
	switch status {
	case "",
		string(domain.OrderStatusPendingPayment),
		string(domain.OrderStatusPaid),
		string(domain.OrderStatusProcessing),
		string(domain.OrderStatusOnDelivery),
		string(domain.OrderStatusCompleted),
		string(domain.OrderStatusCancelled):
		return true
	default:
		return false
	}
}

func (h *OrderHandler) EstimateCheckout(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req domain.CheckoutEstimateRequest
	if len(strings.TrimSpace(string(c.Body()))) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
		}
		if err := h.validate.Struct(&req); err != nil {
			return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
		}
	}

	estimate, err := h.orderUseCase.EstimateCheckout(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyCart) {
			return response.Error(c, fiber.StatusBadRequest, "Keranjang belanja masih kosong")
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Estimasi biaya pengiriman berhasil", estimate)
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req domain.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
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
		return response.Internal(c, err)
	}

	return response.Created(c, "Pesanan berhasil dibuat", order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	page, limit := parsePagination(c)
	status := strings.TrimSpace(c.Query("status", ""))
	if !isValidOrderStatusFilter(status) {
		return response.Error(c, fiber.StatusBadRequest, "Invalid order status filter")
	}

	filter := domain.OrderFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := h.orderUseCase.GetUserOrders(c.Context(), userID, filter)
	if err != nil {
		return response.Internal(c, err)
	}

	return response.Paginated(c, "Orders retrieved successfully", orders, page, limit, total)
}

func (h *OrderHandler) GetOrderByCode(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	code := strings.TrimSpace(c.Params("code"))
	if code == "" {
		code = strings.TrimSpace(strings.TrimPrefix(c.Params("*"), "/"))
	}
	if code == "" || len(code) > 60 {
		return response.Error(c, fiber.StatusBadRequest, "Kode pesanan tidak valid")
	}

	order, err := h.orderUseCase.GetOrderByCode(c.Context(), userID, code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Pesanan tidak ditemukan")
		}
		if errors.Is(err, domain.ErrForbidden) {
			seclog.RecordAccessDenied(c.IP(), userID, c.Method(), c.Path(), "order_ownership_violation")
			return response.Error(c, fiber.StatusForbidden, "Akses pesanan ditolak")
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Detail pesanan berhasil diambil", order)
}
