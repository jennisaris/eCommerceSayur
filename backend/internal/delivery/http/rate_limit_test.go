package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestIPRateLimiterReturns429WhenLimitExceeded(t *testing.T) {
	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "0.0.0.0/0"},
		ProxyHeader:             "X-Real-IP",
	})
	app.Post("/limited", newIPRateLimiter(3, 1*time.Minute, "Too many requests"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/limited", nil)
		req.Header.Set("X-Real-IP", "203.0.113.10")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("attempt %d failed: %v", i, err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("attempt %d: expected 200, got %d", i, resp.StatusCode)
		}
	}

	// 4th request from the same IP must be rejected with 429
	blockedReq := httptest.NewRequest(http.MethodPost, "/limited", nil)
	blockedReq.Header.Set("X-Real-IP", "203.0.113.10")
	resp, err := app.Test(blockedReq)
	if err != nil {
		t.Fatalf("4th attempt failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", resp.StatusCode)
	}

	// Request from a different IP must still succeed
	otherReq := httptest.NewRequest(http.MethodPost, "/limited", nil)
	otherReq.Header.Set("X-Real-IP", "203.0.113.11")
	resp, err = app.Test(otherReq)
	if err != nil {
		t.Fatalf("other IP request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 for separate IP, got %d", resp.StatusCode)
	}
}

func TestUserOrIPRateLimiterScopesPerUser(t *testing.T) {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		if c.Get("X-Test-User") == "1" {
			c.Locals("user_id", int64(1))
		} else if c.Get("X-Test-User") == "2" {
			c.Locals("user_id", int64(2))
		}
		return c.Next()
	})
	app.Post("/pay", newUserOrIPRateLimiter(2, 1*time.Minute, "Too many payment attempts"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	for i := 1; i <= 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/pay", nil)
		req.Header.Set("X-Test-User", "1")
		resp, err := app.Test(req)
		if err != nil || resp.StatusCode != fiber.StatusOK {
			t.Fatalf("user 1 attempt %d: expected 200, got %d err=%v", i, resp.StatusCode, err)
		}
	}

	// 3rd attempt by user 1 -> 429
	req := httptest.NewRequest(http.MethodPost, "/pay", nil)
	req.Header.Set("X-Test-User", "1")
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != fiber.StatusTooManyRequests {
		t.Fatalf("expected 429 for user 1, got %d err=%v", resp.StatusCode, err)
	}

	// User 2 on the same IP still succeeds -> 200
	req2 := httptest.NewRequest(http.MethodPost, "/pay", nil)
	req2.Header.Set("X-Test-User", "2")
	resp, err = app.Test(req2)
	if err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 for user 2, got %d err=%v", resp.StatusCode, err)
	}
}
