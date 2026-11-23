package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

func parsePagination(c *fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	return page, limit
}
