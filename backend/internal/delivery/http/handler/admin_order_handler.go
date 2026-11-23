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

type AdminOrderHandler struct {
	orderUseCase domain.OrderUseCase
	validate     *validator.Validate
}

func NewAdminOrderHandler(orderUseCase domain.OrderUseCase) *AdminOrderHandler {
	return &AdminOrderHandler{
		orderUseCase: orderUseCase,
		validate:     pkgvalidate.New(),
	}
}

func (h *AdminOrderHandler) GetAllOrders(c *fiber.Ctx) error {
	page, limit := parsePagination(c)
	status := strings.TrimSpace(c.Query("status", ""))
	if !isValidOrderStatusFilter(status) {
		return response.Error(c, fiber.StatusBadRequest, "Invalid order status filter")
	}
	search := strings.TrimSpace(c.Query("search", ""))
	if len(search) > 100 {
		search = search[:100]
	}

	filter := domain.OrderFilter{
		Status: status,
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := h.orderUseCase.AdminGetAllOrders(c.Context(), filter)
	if err != nil {
		return response.Internal(c, err)
	}

	return response.Paginated(c, "All orders retrieved successfully", orders, page, limit, total)
}

func (h *AdminOrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return response.Error(c, fiber.StatusBadRequest, "Invalid order ID")
	}

	var req domain.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	order, err := h.orderUseCase.AdminUpdateOrderStatus(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Order not found")
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Order status updated successfully", order)
}
