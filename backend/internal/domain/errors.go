package domain

import "errors"

var (
	ErrNotFound            = errors.New("record not found")
	ErrUnauthorized        = errors.New("unauthorized access")
	ErrForbidden           = errors.New("forbidden: insufficient permissions")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailAlreadyExists  = errors.New("email is already registered")
	ErrUserNotVerified     = errors.New("user account is not verified")
	ErrInsufficientStock   = errors.New("insufficient product stock")
	ErrProductInactive     = errors.New("product is currently inactive")
	ErrInvalidOrderState   = errors.New("invalid order status transition")
	ErrInvalidSignature    = errors.New("invalid midtrans webhook signature")
	ErrPaymentAlreadyPaid  = errors.New("payment has already been processed")
	ErrEmptyCart           = errors.New("cart is empty")
	ErrInvalidAddress      = errors.New("shipping address and recipient details are mandatory")
)
