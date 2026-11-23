package handler

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"
	pkgvalidate "ecommerce-backend/internal/pkg/validate"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUseCase  domain.AuthUseCase
	validate     *validator.Validate
	cookieConfig authCookieConfig
}

func NewAuthHandler(authUseCase domain.AuthUseCase, secure bool, durations ...time.Duration) *AuthHandler {
	accessTTL := 15 * time.Minute
	refreshTTL := 30 * 24 * time.Hour
	if len(durations) > 0 && durations[0] > 0 {
		accessTTL = durations[0]
	}
	if len(durations) > 1 && durations[1] > 0 {
		refreshTTL = durations[1]
	}
	return &AuthHandler{
		authUseCase: authUseCase,
		validate:    pkgvalidate.New(),
		cookieConfig: authCookieConfig{
			secure:        secure,
			accessMaxAge:  cookieMaxAge(accessTTL),
			refreshMaxAge: cookieMaxAge(refreshTTL),
		},
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	res, err := h.authUseCase.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			seclog.Audit(seclog.Event{Event: "auth_register", Outcome: "failure", IP: c.IP(), Email: req.Email, Method: c.Method(), Path: c.Path(), Reason: "email_exists"})
			return response.Error(c, fiber.StatusConflict, err.Error())
		}
		if errors.Is(err, domain.ErrWeakPassword) {
			seclog.Audit(seclog.Event{Event: "auth_register", Outcome: "failure", IP: c.IP(), Email: req.Email, Method: c.Method(), Path: c.Path(), Reason: "weak_password"})
			return response.Error(c, fiber.StatusBadRequest, "Password must be 12-72 characters and include uppercase, lowercase, number, and symbol")
		}
		return response.Internal(c, err)
	}

	var registeredUserID int64
	if res.User != nil {
		registeredUserID = res.User.ID
	}
	seclog.Audit(seclog.Event{Event: "auth_register", Outcome: "success", IP: c.IP(), UserID: registeredUserID, Email: req.Email, Method: c.Method(), Path: c.Path()})

	if res.Token != "" {
		h.setAuthCookies(c, res)
		return response.Created(c, "Registration successful", res)
	}
	return response.Created(c, "Registration successful. Please verify your email", res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	res, err := h.authUseCase.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			seclog.RecordLoginFailure(c.IP(), req.Email, c.Path(), "invalid_credentials")
			return response.Error(c, fiber.StatusUnauthorized, "Invalid email or password")
		}
		if errors.Is(err, domain.ErrUserNotVerified) {
			seclog.RecordLoginFailure(c.IP(), req.Email, c.Path(), "email_not_verified")
			return response.Error(c, fiber.StatusForbidden, "Please verify your email before logging in")
		}
		return response.Internal(c, err)
	}

	var loggedInUserID int64
	if res.User != nil {
		loggedInUserID = res.User.ID
	}
	seclog.Audit(seclog.Event{Event: "auth_login", Outcome: "success", IP: c.IP(), UserID: loggedInUserID, Email: req.Email, Method: c.Method(), Path: c.Path()})

	h.setAuthCookies(c, res)
	return response.Success(c, "Login successful", res)
}

func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token" validate:"required,min=16,max=256"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Verification token is required")
	}
	if err := h.authUseCase.VerifyEmail(c.Context(), req.Token); err != nil {
		seclog.Audit(seclog.Event{Event: "auth_verify_email", Outcome: "failure", IP: c.IP(), Method: c.Method(), Path: c.Path(), Reason: "invalid_or_expired_token"})
		return response.Error(c, fiber.StatusBadRequest, "Invalid or expired verification token")
	}
	seclog.Audit(seclog.Event{Event: "auth_verify_email", Outcome: "success", IP: c.IP(), Method: c.Method(), Path: c.Path()})
	return response.Success(c, "Email verified successfully", nil)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	res, err := h.authUseCase.Refresh(c.Context(), c.Cookies(refreshTokenCookie))
	if err != nil {
		seclog.Audit(seclog.Event{Event: "auth_refresh", Outcome: "failure", IP: c.IP(), Method: c.Method(), Path: c.Path(), Reason: "invalid_or_expired_session"})
		h.clearAuthCookies(c)
		return response.Error(c, fiber.StatusUnauthorized, "Session expired. Please login again")
	}

	h.setAuthCookies(c, res)
	return response.Success(c, "Session refreshed", res)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.authUseCase.Logout(c.Context(), c.Cookies(refreshTokenCookie)); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Unable to logout")
	}
	h.clearAuthCookies(c)
	seclog.Audit(seclog.Event{Event: "auth_logout", Outcome: "success", IP: c.IP(), Method: c.Method(), Path: c.Path()})
	return response.Success(c, "Logout successful", nil)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	user, err := h.authUseCase.GetProfile(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "User not found")
	}

	return response.Success(c, "Profile retrieved successfully", user)
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req domain.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, pkgvalidate.FormatError(err))
	}

	user, emailChanged, err := h.authUseCase.UpdateProfile(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			seclog.Audit(seclog.Event{Event: "profile_update", Outcome: "failure", IP: c.IP(), UserID: userID, Email: req.Email, Method: c.Method(), Path: c.Path(), Reason: "invalid_current_password"})
			return response.Error(c, fiber.StatusUnauthorized, "Kata sandi saat ini tidak valid untuk mengubah alamat email")
		}
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return response.Error(c, fiber.StatusConflict, "Alamat email sudah digunakan oleh akun lain")
		}
		if errors.Is(err, domain.ErrNotFound) {
			return response.Error(c, fiber.StatusNotFound, "User not found")
		}
		return response.Internal(c, err)
	}

	eventName := "profile_updated"
	if emailChanged {
		eventName = "profile_email_changed"
	}
	seclog.Audit(seclog.Event{Event: eventName, Outcome: "success", IP: c.IP(), UserID: userID, Email: user.Email, Method: c.Method(), Path: c.Path()})
	return response.Success(c, "Profil berhasil diperbarui", user)
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req domain.ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Masukkan alamat email yang valid")
	}

	if err := h.authUseCase.RequestPasswordReset(c.Context(), req.Email); err != nil {
		return response.Internal(c, err)
	}

	seclog.Audit(seclog.Event{Event: "password_reset_requested", Outcome: "success", IP: c.IP(), Email: req.Email, Method: c.Method(), Path: c.Path()})
	return response.Success(c, "Jika email terdaftar, tautan reset password telah dikirim.", nil)
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req domain.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if err := h.validate.Struct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Token tidak valid atau password baru tidak memenuhi syarat (minimal 12 karakter dengan huruf besar, huruf kecil, angka, dan simbol)")
	}

	if err := h.authUseCase.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		if errors.Is(err, domain.ErrWeakPassword) {
			seclog.Audit(seclog.Event{Event: "password_reset_completed", Outcome: "failure", IP: c.IP(), Method: c.Method(), Path: c.Path(), Reason: "weak_password"})
			return response.Error(c, fiber.StatusBadRequest, "Password harus 12-72 karakter serta mengandung huruf besar, huruf kecil, angka, dan simbol")
		}
		if errors.Is(err, domain.ErrInvalidResetToken) {
			seclog.Audit(seclog.Event{Event: "password_reset_completed", Outcome: "failure", IP: c.IP(), Method: c.Method(), Path: c.Path(), Reason: "invalid_or_expired_token"})
			return response.Error(c, fiber.StatusBadRequest, "Tautan reset password tidak valid atau sudah kedaluwarsa")
		}
		return response.Internal(c, err)
	}

	h.clearAuthCookies(c)
	seclog.Audit(seclog.Event{Event: "password_reset_completed", Outcome: "success", IP: c.IP(), Method: c.Method(), Path: c.Path()})
	return response.Success(c, "Password berhasil diperbarui. Silakan login kembali dengan password baru Anda.", nil)
}
