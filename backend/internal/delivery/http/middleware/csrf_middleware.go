package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"
	"github.com/gofiber/fiber/v2"
)

const csrfCookieName = "csrf_token"

func newCSRFToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func CSRFToken(c *fiber.Ctx) (string, error) {
	if token := c.Cookies(csrfCookieName); token != "" {
		return token, nil
	}

	token, err := newCSRFToken()
	if err != nil {
		return "", err
	}
	c.Cookie(&fiber.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   c.Locals("cookie_secure") == true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	return token, nil
}

func parseAllowedOrigins(allowedOriginsCSV ...string) map[string]struct{} {
	allowed := make(map[string]struct{})
	for _, rawList := range allowedOriginsCSV {
		for _, item := range strings.Split(rawList, ",") {
			if normalized := normalizeOrigin(item); normalized != "" {
				allowed[normalized] = struct{}{}
			}
		}
	}
	return allowed
}

func normalizeOrigin(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return strings.ToLower(parsed.Scheme) + "://" + strings.ToLower(parsed.Host)
}

func CSRFProtection(secure bool, allowedOriginsCSV ...string) fiber.Handler {
	allowedOrigins := parseAllowedOrigins(allowedOriginsCSV...)

	return func(c *fiber.Ctx) error {
		c.Locals("cookie_secure", secure)
		if _, err := CSRFToken(c); err != nil {
			return response.Error(c, fiber.StatusInternalServerError, "Unable to initialize security token")
		}

		if c.Method() == http.MethodGet || c.Method() == http.MethodHead || c.Method() == http.MethodOptions || c.Path() == "/api/v1/payments/midtrans-callback" {
			return c.Next()
		}

		origin := strings.TrimSpace(c.Get("Origin"))
		referer := strings.TrimSpace(c.Get("Referer"))
		originTrusted := false

		if len(allowedOrigins) > 0 {
			if origin != "" {
				if _, ok := allowedOrigins[normalizeOrigin(origin)]; !ok {
					seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "csrf_invalid_origin")
					return response.Error(c, fiber.StatusForbidden, "Invalid request origin")
				}
				originTrusted = true
			} else if referer != "" {
				if _, ok := allowedOrigins[normalizeOrigin(referer)]; !ok {
					seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "csrf_invalid_referer")
					return response.Error(c, fiber.StatusForbidden, "Invalid request origin")
				}
				originTrusted = true
			}
		}

		if strings.EqualFold(strings.TrimSpace(c.Get("Sec-Fetch-Site")), "cross-site") && !originTrusted {
			seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "csrf_cross_site_blocked")
			return response.Error(c, fiber.StatusForbidden, "Cross-site request blocked")
		}

		cookieToken := c.Cookies(csrfCookieName)
		headerToken := c.Get("X-CSRF-Token")
		if cookieToken == "" || headerToken == "" || len(cookieToken) != len(headerToken) || subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "csrf_token_mismatch")
			return response.Error(c, fiber.StatusForbidden, "Invalid CSRF token")
		}
		return c.Next()
	}
}
