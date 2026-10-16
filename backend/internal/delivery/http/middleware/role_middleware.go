package middleware

import (
	"ecommerce-backend/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rolesVal := c.Locals("roles")
		if rolesVal == nil {
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: no roles found")
		}

		roles, ok := rolesVal.([]string)
		if !ok {
			return response.Error(c, fiber.StatusForbidden, "Access forbidden: invalid role format")
		}

		for _, r := range roles {
			if r == requiredRole {
				return c.Next()
			}
		}

		return response.Error(c, fiber.StatusForbidden, "Access forbidden: required role '"+requiredRole+"' not met")
	}
}
