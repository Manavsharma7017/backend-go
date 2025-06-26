package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func InterviewSessionRoute(c *fiber.App) {
	sessionGroup := c.Group("/api/sessions", middleware.UserAuthMiddleware)
	sessionGroup.Post("/create", handlers.StartInterviewSession)
	sessionGroup.Get("/getall", handlers.GetAllInterviewSession)
	sessionGroup.Get("/:id", handlers.GetInterviewSessionById)
	sessionGroup.Patch("/:id/complete", handlers.EndInterviewSession)
}
