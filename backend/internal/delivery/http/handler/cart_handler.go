package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartUseCase domain.CartUseCase
	validate    *validator.Validate
}

func NewCartHandler(cartUseCase domain.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUseCase: cartUseCase,
		validate:    validator.New(),
	}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	summary, err := h.cartUseCase.GetCart(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return response.Success(c, "Cart retrieved successfully", summary)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	var req domain.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	summary, err := h.cartUseCase.AddToCart(c.Context(), userID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "Item added to cart", summary)
}

func (h *CartHandler) UpdateItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	productID, err := strconv.ParseInt(c.Params("product_id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	var req domain.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	summary, err := h.cartUseCase.UpdateItem(c.Context(), userID, productID, &req)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "Cart item updated", summary)
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	productID, err := strconv.ParseInt(c.Params("product_id"), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid product ID")
	}

	summary, err := h.cartUseCase.RemoveItem(c.Context(), userID, productID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Item removed from cart", summary)
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	if err := h.cartUseCase.ClearCart(c.Context(), userID); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, "Cart cleared successfully", nil)
}
