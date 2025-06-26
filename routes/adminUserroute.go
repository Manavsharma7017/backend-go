package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminUserRoutes(c *fiber.App) {
	adminGroup := c.Group("/api/admin")
	adminGroup.Post("/signup", handlers.CreateAdminUserHandler)
	adminGroup.Post("/login", handlers.LoginAdminUserHandler)
	protected := adminGroup.Group("/", middleware.AdminAuthMiddleware)
	protected.Get("/profile", handlers.GetAllAdminUsersHandler)
	protected.Get("/dashboard", handlers.GetAdminDashboardHandler)
}
