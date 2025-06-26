package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateDomainHandler creates a new domain
func CreateDomainHandler(c *fiber.Ctx) error {
	var domain models.Domain
	if err := c.BodyParser(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := utils.ValidateStruct(domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	_, err := services.CreateDomain(&domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Domain created successfully",
		"domain":  domain,
	})
}

// GetDomainHandler retrieves a domain by ID
func GetDomainHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	domain, err := services.GetDomainByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Domain not found",
		})
	}
	return c.JSON(fiber.Map{
		"domain": domain,
	})
}

// GetAllDomainsHandler retrieves all domains
func GetAllDomainsHandler(c *fiber.Ctx) error {
	domains, err := services.GetAllDomains()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"domains": domains,
	})
}

// UpdateDomainHandler updates a domain by ID
func UpdateDomainHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var domain models.Domain
	if err := c.BodyParser(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updatedDomain, err := services.UpdateDomain(id, &domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Domain updated successfully",
		"domain":  updatedDomain,
	})
}

// DeleteDomainHandler deletes a domain by ID
func DeleteDomainHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := services.DeleteDomain(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Domain deleted successfully",
	})
}
