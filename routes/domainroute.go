package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func DomainRoutes(c *fiber.App) {
	domainGroup := c.Group("/api/domain", middleware.UserAuthMiddleware)
	domainAdminGroup := c.Group("/api/admin/domain", middleware.AdminAuthMiddleware)
	domainAdminGroup.Post("/create", handlers.CreateDomainHandler)
	domainAdminGroup.Put("/update/:id", handlers.UpdateDomainHandler)
	domainAdminGroup.Delete("/delete/:id", handlers.DeleteDomainHandler)
	domainAdminGroup.Get("/get/:id", handlers.GetDomainHandler)
	domainAdminGroup.Get("/getall", handlers.GetAllDomainsHandler)

	domainGroup.Get("/get/:id", handlers.GetDomainHandler)
	domainGroup.Get("/getall", handlers.GetAllDomainsHandler)

}
