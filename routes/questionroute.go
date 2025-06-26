package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func QuestionRoutes(c *fiber.App) {
	questiongroup := c.Group("/api/questions", middleware.UserAuthMiddleware)
	questionAdmingroup := c.Group("/api/admin/questions", middleware.AdminAuthMiddleware)
	questionAdmingroup.Post("/", handlers.CreateQuestionHandler)
	questionAdmingroup.Get("/", handlers.GetAllQuestionsHandler)
	questionAdmingroup.Put("/:id", handlers.UpdateQuestionHandler)
	questionAdmingroup.Delete("/:id", handlers.DeleteQuestionHandler)
	questiongroup.Get("/", handlers.GetQuestionFilterHandler)
	questiongroup.Get("/:id", handlers.GetQuestionHandler)
	questiongroup.Get("/find/all", handlers.GetAllQuestionsHandler)

}
