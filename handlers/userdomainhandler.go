package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func CreateUserDomainHandler(c *fiber.Ctx) error {
	var userDomain models.UserDomain

	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	userDomain.UserID = userID

	if err := c.BodyParser(&userDomain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	if err := utils.ValidateStruct(&userDomain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// ✅ Get actual created or existing domain
	savedUserDomain, err := services.CreateUserDomain(&userDomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User domain created or reused successfully",
		"user_domain": fiber.Map{
			"id":         savedUserDomain.ID,
			"domain_id":  savedUserDomain.DomainID,
			"user_id":    savedUserDomain.UserID,
			"created_at": savedUserDomain.CreatedAt,
		},
	})
}

func GetAllUserDomainsHandler(c *fiber.Ctx) error {
	// Handler logic for getting all user domains
	var userdomains []models.UserDomain
	var user_id string = c.Locals("user_id").(string)
	if user_id == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	_, err := services.GetAllUserDomains(&userdomains, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	domainchan := make(chan models.UserDomain)
	var sanitizedata []fiber.Map
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userd := range domainchan {
				sanitized := fiber.Map{
					"id":        userd.ID,
					"domain_id": userd.DomainID,
					"user_id":   userd.UserID,
					"domain": fiber.Map{
						"id":   userd.Domain.ID,
						"name": userd.Domain.Name,
					},

					"created_at": userd.CreatedAt,
				}
				mu.Lock()
				sanitizedata = append(sanitizedata, sanitized)
				mu.Unlock()
			}
		}()
	}
	for _, userd := range userdomains {
		domainchan <- userd
	}
	close(domainchan)
	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "User domains fetched successfully",
		"user_domains": sanitizedata,
	})

}
func GetUserDomainByIDHandler(c *fiber.Ctx) error {
	// Handler logic for getting a user domain by ID
	id := c.Params("id")

	var userDomain models.UserDomain
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User domain ID is required",
		})
	}
	_, err := services.GetUserDomainByID(&userDomain, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "User domain fetched successfully",
		"user_domain": fiber.Map{"id": userDomain.ID, "domain_id": userDomain.DomainID, "user_id": userDomain.UserID, "created_at": userDomain.CreatedAt},
	})

}

func DeleteUserDomainHandler(c *fiber.Ctx) error {

	id := c.Params("id")

	var userDomain models.UserDomain
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User domain ID is required",
		})
	}
	_, err := services.DeletUserDomainByID(&userDomain, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User domain deleted successfully",
	})
}
