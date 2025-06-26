package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserDomainRoute(c *fiber.App) {
	userDomainGroup := c.Group("/api/user-domains", middleware.UserAuthMiddleware)

	userDomainGroup.Post("/", handlers.CreateUserDomainHandler)
	userDomainGroup.Get("/", handlers.GetAllUserDomainsHandler)
	userDomainGroup.Get("/:id", handlers.GetUserDomainByIDHandler)
	userDomainGroup.Delete("/:id", handlers.DeleteUserDomainHandler)
}
