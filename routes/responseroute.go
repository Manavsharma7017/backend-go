package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func ResponseRoute(c *fiber.App) {
	responceGroup := c.Group("/api/responses", middleware.UserAuthMiddleware)
	responceSessionGroup := c.Group("/api/session", middleware.UserAuthMiddleware)
	responceGroup.Post("/", handlers.CreateResponseHandler)
	responceGroup.Get("/:id", handlers.GetResponseHandler)
	responceSessionGroup.Get("/:id/responses", handlers.GetSessionResponseHandler)

}
