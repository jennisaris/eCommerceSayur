package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"
	"errors"
	"log"
	"strings"

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
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}
	orderCode := c.Params("order_code")
	if orderCode == "" {
		wildcard := c.Params("*")
		orderCode = strings.TrimSuffix(wildcard, "/snap")
		orderCode = strings.TrimPrefix(orderCode, "/")
	}
	if orderCode == "" {
		orderCode = c.Query("order_code")
	}
	if orderCode == "" {
		var req struct {
			OrderCode string `json:"order_code"`
		}
		_ = c.BodyParser(&req)
		orderCode = req.OrderCode
	}

	orderCode = strings.TrimSpace(orderCode)
	if orderCode == "" || len(orderCode) > 60 {
		return response.Error(c, fiber.StatusBadRequest, "Kode pesanan tidak valid")
	}

	res, err := h.paymentUseCase.CreateSnapTransaction(c.Context(), orderCode, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Pesanan tidak ditemukan")
		}
		if errors.Is(err, domain.ErrForbidden) {
			seclog.RecordAccessDenied(c.IP(), userID, c.Method(), c.Path(), "payment_ownership_violation")
			return response.Error(c, fiber.StatusForbidden, "Akses pesanan ditolak")
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return response.Success(c, "Snap token berhasil dibuat", res)
}

func (h *PaymentHandler) MidtransCallback(c *fiber.Ctx) error {
	rawPayload := c.Body()

	log.Printf("[MIDTRANS WEBHOOK RECEIVED] payload_bytes=%d", len(rawPayload))

	err := h.paymentUseCase.HandleMidtransCallback(c.Context(), rawPayload)
	if err != nil {
		log.Printf("[MIDTRANS WEBHOOK ERROR] %v", err)
		if errors.Is(err, domain.ErrInvalidSignature) {
			seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "invalid_midtrans_webhook_signature")
			return response.Error(c, fiber.StatusUnauthorized, "Invalid webhook signature")
		}
		return response.Error(c, fiber.StatusBadRequest, "Invalid webhook payload")
	}

	return response.Success(c, "Midtrans webhook handled successfully", nil)
}

func (h *PaymentHandler) SyncPaymentStatus(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	orderCode := strings.TrimSpace(c.Params("order_code"))
	if orderCode == "" {
		var req struct {
			OrderCode string `json:"order_code"`
		}
		_ = c.BodyParser(&req)
		orderCode = strings.TrimSpace(req.OrderCode)
	}
	if orderCode == "" || len(orderCode) > 60 {
		return response.Error(c, fiber.StatusBadRequest, "Kode pesanan tidak valid")
	}

	order, err := h.paymentUseCase.SyncOrderPaymentStatus(c.Context(), orderCode, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "Pesanan tidak ditemukan")
		}
		if errors.Is(err, domain.ErrForbidden) {
			seclog.RecordAccessDenied(c.IP(), userID, c.Method(), c.Path(), "payment_sync_ownership_violation")
			return response.Error(c, fiber.StatusForbidden, "Akses pesanan ditolak")
		}
		return response.Internal(c, err)
	}

	return response.Success(c, "Status pembayaran berhasil disinkronkan", order)
}
