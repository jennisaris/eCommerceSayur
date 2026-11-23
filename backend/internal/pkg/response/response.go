package response

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	CurrentPage int   `json:"current_page,omitempty"`
	PerPage     int   `json:"per_page,omitempty"`
	TotalItems  int64 `json:"total_items,omitempty"`
	TotalPages  int   `json:"total_pages,omitempty"`
}

func JSON(c *fiber.Ctx, statusCode int, success bool, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusOK, true, message, data)
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
	return JSON(c, fiber.StatusCreated, true, message, data)
}

func Error(c *fiber.Ctx, statusCode int, message string) error {
	return JSON(c, statusCode, false, message, nil)
}

func Internal(c *fiber.Ctx, err error) error {
	log.Printf("[INTERNAL_ERROR] method=%s path=%s error=%v", c.Method(), c.Path(), err)
	return Error(c, fiber.StatusInternalServerError, "Internal server error")
}

func Paginated(c *fiber.Ctx, message string, data interface{}, page, limit int, totalItems int64) error {
	totalPages := 0
	if limit > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			CurrentPage: page,
			PerPage:     limit,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})
}
