package middleware

import (
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/response"
	"ecommerce-backend/internal/pkg/seclog"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(requiredRole string, userRepo domain.UserRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(int64)
		if !ok || userID == 0 {
			seclog.RecordAccessDenied(c.IP(), 0, c.Method(), c.Path(), "unauthenticated_admin_access")
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		user, err := userRepo.GetByID(c.Context(), userID)
		if err != nil || user == nil || !user.HasRole(requiredRole) {
			seclog.RecordAccessDenied(c.IP(), userID, c.Method(), c.Path(), "missing_required_role:"+requiredRole)
			return response.Error(c, fiber.StatusForbidden, "Access forbidden")
		}

		return c.Next()
	}
}
