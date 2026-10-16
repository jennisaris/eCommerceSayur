package middleware

import (
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtSvc *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token format. Expected 'Bearer <token>'")
		}

		claims, err := jwtSvc.ValidateToken(parts[1])
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)

		return c.Next()
	}
}
