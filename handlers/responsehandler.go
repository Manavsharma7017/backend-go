package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func CreateResponseHandler(c *fiber.Ctx) error {
	var responce models.Response
	if err := c.BodyParser(&responce); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse request body",
		})
	}
	if err := utils.ValidateStruct(responce); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "failed to validate request body",
			"details": err.Error(),
		})
	}
	_, err := services.CreateResponse(&responce)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create responce",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "responce created",
		"responce": fiber.Map{
			"responceID":  responce.ResponseID,
			"sessionID":   responce.SessionID,
			"questionID":  responce.UserQuestionID,
			"answer":      responce.Answer,
			"submittedAt": responce.SubmittedAt,
		},
	})

}
func GetResponseHandler(c *fiber.Ctx) error {
	responceID := c.Params("id")
	if responceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "responce id is required",
		})
	}
	responce, err := services.GetResponse(responceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get responce",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "responce fetched",
		"responce": fiber.Map{
			"responceID":  responce.ResponseID,
			"sessionID":   responce.SessionID,
			"questionID":  responce.UserQuestionID,
			"answer":      responce.Answer,
			"Feedback":    responce.Feedback,
			"submittedAt": responce.SubmittedAt,
		},
	})

}
func GetSessionResponseHandler(c *fiber.Ctx) error {
	sessionID := c.Params("id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "session id is required",
		})
	}
	responce, err := services.GetSessionResponse(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to get responce",
			"details": err.Error(),
		})
	}
	responceChan := make((chan models.Response))
	var sanitizeddata []fiber.Map
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for responce := range responceChan {
				sanitized := fiber.Map{
					"responceID":  responce.ResponseID,
					"sessionID":   responce.SessionID,
					"questionID":  responce.UserQuestionID,
					"answer":      responce.Answer,
					"submittedAt": responce.SubmittedAt,
				}
				mu.Lock()
				sanitizeddata = append(sanitizeddata, sanitized)
				mu.Unlock()
			}
		}()
	}
	for _, responce := range responce {
		responceChan <- responce
	}
	close(responceChan)
	wg.Wait()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "responce fetched",
		"responce": sanitizeddata,
	})
}
