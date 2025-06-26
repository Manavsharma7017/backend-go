package middleware

import (
	"backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AdminAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	tokenstring := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseAdminJWT(tokenstring)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}
	role, ok := claims["role"].(string)
	if !ok || role != "ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: admin only",
		})
	}
	c.Locals("user_id", claims["user_id"])
	c.Locals("role", role)

	return c.Next()
}
