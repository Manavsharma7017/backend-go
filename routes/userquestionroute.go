package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserQuestionRoute(c *fiber.App) {
	userQuestionGroup := c.Group("/api/user-questions", middleware.UserAuthMiddleware)

	userQuestionGroup.Post("/", handlers.CreateUserQuestionHandler)
	userQuestionGroup.Get("/", handlers.GetAllUserQuestionsHandler)
	userQuestionGroup.Get("/:id", handlers.GetUserQuestionByIDHandler)
	userQuestionGroup.Delete("/:id", handlers.DeleteUserQuestionHandler)
}
