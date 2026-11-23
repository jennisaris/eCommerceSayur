package handler

import "github.com/gofiber/fiber/v2"

func authenticatedUserID(c *fiber.Ctx) (int64, bool) {
	userID, ok := c.Locals("user_id").(int64)
	if !ok || userID <= 0 {
		return 0, false
	}
	return userID, true
}
