package handlers

import (
	"backend/models"
	"backend/services"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func CreateUserQuestionHandler(c *fiber.Ctx) error {
	// Handler logic for creating a user question
	var userQuestion models.UserQuestion
	var user_id string = c.Locals("user_id").(string)

	// Parse the request body into the userQuestion model
	if err := c.BodyParser(&userQuestion); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Create user question via service layer
	_, err := services.CreateUserQuestion(&userQuestion, user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user question",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User question created successfully",
		"user_question": fiber.Map{
			"id":             userQuestion.ID,
			"questionNumber": userQuestion.Questionnumber,
			"sessionID":      userQuestion.SessionID,
			"Text":           userQuestion.Text,
			"Explanation":    userQuestion.Explanation,
			"Difficulty":     userQuestion.Difficulty,
		}})
}

func GetAllUserQuestionsHandler(c *fiber.Ctx) error {
	// Handler logic for getting all user questions
	userQuestions, err := services.GetAllUserQuestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user questions",
		})
	}
	userQuestionschan := make(chan models.UserQuestion)
	var wg sync.WaitGroup
	var sanitizedata []fiber.Map
	var mu sync.Mutex
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, userQuestion := range userQuestions {
				sanatied := fiber.Map{
					"id":             userQuestion.ID,
					"questionNumber": userQuestion.Questionnumber,
					"sessionID":      userQuestion.SessionID,

					"createdAt": userQuestion.CreatedAt,
					"updatedAt": userQuestion.UpdatedAt,

					"userID": userQuestion.UserID,
				}
				mu.Lock()
				sanitizedata = append(sanitizedata, sanatied)
				mu.Unlock()
			}
		}()
	}
	for _, userQuestion := range userQuestions {
		userQuestionschan <- userQuestion
	}
	close(userQuestionschan)
	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(sanitizedata)
}

func GetUserQuestionByIDHandler(c *fiber.Ctx) error {

	id := c.Params("id")
	userQuestion, err := services.GetUserQuestionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User question not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userQuestion)
}
func DeleteUserQuestionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	err := services.DeleteUserQuestion(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user question",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User question deleted successfully",
	})
}
