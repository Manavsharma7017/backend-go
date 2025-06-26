package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(c *fiber.App) {

	userGroup := c.Group("/api/user")
	userGroup.Post("/signup", handlers.CreateUserHandler)
	userGroup.Post("/login", handlers.LoginUserHandler)
	protectedGroup := c.Group("/api/users", middleware.UserAuthMiddleware)
	protectedGroup.Get("/", handlers.GetUserHandler)
}
