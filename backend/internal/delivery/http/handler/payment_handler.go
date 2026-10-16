package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentUseCase domain.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase domain.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

func (h *PaymentHandler) CreateSnapToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	orderCode := c.Params("order_code")

	res, err := h.paymentUseCase.CreateSnapTransaction(c.Context(), orderCode, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Pesanan tidak ditemukan")
		}
		if errors.Is(err, domain.ErrForbidden) {
			return response.Error(c, fiber.StatusForbidden, "Akses pesanan ditolak")
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "Snap token berhasil dibuat", res)
}

func (h *PaymentHandler) MidtransCallback(c *fiber.Ctx) error {
	rawPayload := c.Body()

	log.Printf("[MIDTRANS WEBHOOK RECEIVED] Payload: %s", string(rawPayload))

	err := h.paymentUseCase.HandleMidtransCallback(c.Context(), rawPayload)
	if err != nil {
		log.Printf("[MIDTRANS WEBHOOK ERROR] %v", err)
		if errors.Is(err, domain.ErrInvalidSignature) {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid webhook signature")
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "Midtrans webhook handled successfully", nil)
}
