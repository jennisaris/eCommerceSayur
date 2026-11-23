package domain

import "errors"

var (
	ErrNotFound           = errors.New("record not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden: insufficient permissions")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidSession     = errors.New("invalid or expired session")
	ErrInvalidResetToken  = errors.New("invalid or expired password reset token")
	ErrWeakPassword       = errors.New("password does not meet security requirements")
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrUserNotVerified    = errors.New("user account is not verified")
	ErrInsufficientStock  = errors.New("insufficient product stock")
	ErrProductInactive    = errors.New("product is currently inactive")
	ErrInvalidOrderState  = errors.New("invalid order status transition")
	ErrInvalidSignature   = errors.New("invalid midtrans webhook signature")
	ErrPaymentAlreadyPaid = errors.New("payment has already been processed")
	ErrEmptyCart          = errors.New("cart is empty")
	ErrInvalidAddress     = errors.New("shipping address and recipient details are mandatory")
)
