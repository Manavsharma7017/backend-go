package middleware

import (
	"backend/utils"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UserAuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := utils.ParseUserJWT(tokenString)
	if err != nil {
		log.Printf("❌ JWT parse failed: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		log.Println("❌ user_id not found in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	c.Locals("user_id", userID)

	return c.Next()
}
