package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func FeedBackRoute(c *fiber.App) {
	feedbackGroup := c.Group("/api/responses", middleware.UserAuthMiddleware)
	feedbackGroup.Post("/feedback", handlers.CreateFeedbackHandler)
	feedbackGroup.Get("/:id/feedback", handlers.GetFeedbackHandler)

}
