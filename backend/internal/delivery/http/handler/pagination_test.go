package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestParsePaginationUsesSafeBounds(t *testing.T) {
	app := fiber.New()
	app.Get("/items", func(c *fiber.Ctx) error {
		page, limit := parsePagination(c)
		return c.JSON(fiber.Map{"page": page, "limit": limit})
	})

	request := httptest.NewRequest("GET", "/items?page=-2&limit=1000", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer response.Body.Close()

	var result struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("failed decoding response: %v", err)
	}
	if result.Page != defaultPage || result.Limit != maxPageSize {
		t.Fatalf("expected page=%d and limit=%d, got page=%d and limit=%d", defaultPage, maxPageSize, result.Page, result.Limit)
	}
}
