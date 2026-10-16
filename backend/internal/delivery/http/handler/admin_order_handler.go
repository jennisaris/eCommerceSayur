package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"errors"
	"strconv"

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
		validate:     validator.New(),
	}
}

func (h *AdminOrderHandler) GetAllOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")
	search := c.Query("search", "")

	filter := domain.OrderFilter{
		Status: status,
		Search: search,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := h.orderUseCase.AdminGetAllOrders(c.Context(), filter)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Paginated(c, "All orders retrieved successfully", orders, page, limit, total)
}

func (h *AdminOrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid order ID")
	}

	var req domain.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	order, err := h.orderUseCase.AdminUpdateOrderStatus(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Order not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Order status updated successfully", order)
}
