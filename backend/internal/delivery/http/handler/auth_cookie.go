package handler

import (
	"ecommerce-backend/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

type authCookieConfig struct {
	secure        bool
	accessMaxAge  int
	refreshMaxAge int
}

func (h *AuthHandler) setAuthCookies(c *fiber.Ctx, auth *domain.AuthResponse) {
	c.Cookie(&fiber.Cookie{
		Name:     accessTokenCookie,
		Value:    auth.Token,
		Path:     "/",
		MaxAge:   h.cookieConfig.accessMaxAge,
		HTTPOnly: true,
		Secure:   h.cookieConfig.secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookie,
		Value:    auth.RefreshToken,
		Path:     "/api/v1/auth",
		MaxAge:   h.cookieConfig.refreshMaxAge,
		HTTPOnly: true,
		Secure:   h.cookieConfig.secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func (h *AuthHandler) clearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cookieConfig.secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cookieConfig.secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func cookieMaxAge(duration time.Duration) int {
	return int(duration / time.Second)
}
