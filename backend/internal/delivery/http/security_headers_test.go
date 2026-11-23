package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeadersMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(newSecurityHeadersMiddleware(true))
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/health", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	checks := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Content-Security-Policy":   "frame-ancestors 'none'",
		"Strict-Transport-Security": "max-age=31536000",
	}
	for header, wantSubstr := range checks {
		got := resp.Header.Get(header)
		if !strings.Contains(got, wantSubstr) {
			t.Errorf("header %s = %q, expected to contain %q", header, got, wantSubstr)
		}
	}
}
