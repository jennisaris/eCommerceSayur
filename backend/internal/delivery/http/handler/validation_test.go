package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestParseProductFilterValidation(t *testing.T) {
	app := fiber.New()
	app.Get("/products", func(c *fiber.Ctx) error {
		_, _, _, err := parseProductFilter(c, false)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})
	app.Get("/admin/products", func(c *fiber.Ctx) error {
		_, _, _, err := parseProductFilter(c, true)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	cases := []struct {
		url      string
		wantCode int
	}{
		{"/products?min_price=10000&max_price=50000&sort=price_asc", fiber.StatusOK},
		{"/products?min_price=-100", fiber.StatusBadRequest},
		{"/products?min_price=50000&max_price=10000", fiber.StatusBadRequest},
		{"/products?category_id=0", fiber.StatusBadRequest},
		{"/products?sort=drop_table", fiber.StatusBadRequest},
		{"/admin/products?status=all", fiber.StatusOK},
		{"/admin/products?status=HACKED", fiber.StatusBadRequest},
	}

	for _, tc := range cases {
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, tc.url, nil))
		if err != nil {
			t.Fatalf("request %s failed: %v", tc.url, err)
		}
		if resp.StatusCode != tc.wantCode {
			t.Errorf("GET %s: expected status %d, got %d", tc.url, tc.wantCode, resp.StatusCode)
		}
	}
}

func TestIsValidOrderStatusFilter(t *testing.T) {
	if !isValidOrderStatusFilter("") || !isValidOrderStatusFilter("PAID") || !isValidOrderStatusFilter("ON_DELIVERY") {
		t.Fatal("expected valid order statuses to pass")
	}
	if isValidOrderStatusFilter("UNKNOWN_STATUS") || isValidOrderStatusFilter("PAID' OR '1'='1") {
		t.Fatal("expected invalid order statuses to be rejected")
	}
}
