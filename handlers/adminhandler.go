package handlers

import (
	"backend/config"
	"backend/models"
	"backend/services"
	"backend/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func CreateAdminUserHandler(c *fiber.Ctx) error {
	var user models.AdminUser
	Adminpassword := c.Get("Admin-Password")
	if Adminpassword != config.GetAdminpassword() {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Forbidden: Invalid Admin Password",
		})
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input data",
		})
	}

	if err := utils.ValidateStruct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"details": err.Error(),
		})
	}

	_, err := services.CreateAdminUserService(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	token, err := utils.GetAdminJWT(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin user created successfully",
		"token":   token,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})

}

func LoginAdminUserHandler(c *fiber.Ctx) error {
	var adminuserlogin models.LoginModel

	if err := c.BodyParser(&adminuserlogin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid input data",
		})
	}
	if err := utils.ValidateStruct(adminuserlogin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Validation failed",
			"details": err.Error(),
		})
	}
	user, err := services.LoginAdminUserService(adminuserlogin)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	token, err := utils.GetAdminJWT(user.ID, string(user.Role))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin user Login successfully",
		"token":   token,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})

}
func GetAdminDashboardHandler(c *fiber.Ctx) error {
	dashboardData, err := services.GetAdminDashboardService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve dashboard data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin Dashboard Data",
		"data":    dashboardData,
	})
}
func GetAllAdminUsersHandler(c *fiber.Ctx) error {
	adminuser, err := services.GetAllAdminUserService()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Admin user not found",
		})

	}
	userChan := make(chan models.AdminUser)
	var sanitizedUsers []fiber.Map
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Start 10 workers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range userChan { // <-- receive users one by one
				sanitized := fiber.Map{
					"id":       user.ID,
					"email":    user.Email,
					"username": user.Username,
					"role":     user.Role,
				}
				mu.Lock()
				sanitizedUsers = append(sanitizedUsers, sanitized)
				mu.Unlock()
			}
		}()
	}

	// Feed users into channel
	for _, user := range adminuser {
		userChan <- user // <-- sending ONE user at a time
	}
	close(userChan) // close after all users are sent

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Get All Admin Users Handler",
		"adminuserdata": sanitizedUsers,
	})
}
