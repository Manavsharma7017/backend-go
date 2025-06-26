package handlers

import (
	"backend/models"
	"backend/services"

	"github.com/gofiber/fiber/v2"
)

func GetFeedbackHandler(c *fiber.Ctx) error {
	var ResponseID = c.Params("id")
	var Feedback models.Feedback
	if ResponseID == "" {
		return c.Status(400).SendString("ID is required")
	}
	_, err := services.GetFeedbackByResponseID(ResponseID, &Feedback)
	if err != nil {
		return c.Status(500).SendString("Error fetching feedback")
	}
	if Feedback.ID == "" {
		return c.Status(404).SendString("Feedback not found")
	}

	return c.Status(fiber.StatusCreated).JSON(Feedback)
}
func CreateFeedbackHandler(c *fiber.Ctx) error {
	var userrequest models.UserQuestionfedback
	var userFeedback models.Feedback
	if err := c.BodyParser(&userrequest); err != nil {
		return c.Status(400).SendString("Invalid request body")
	}
	data, err := services.CreateFeedback(&userFeedback, userrequest)
	if err != nil {
		return c.Status(500).SendString("Error creating feedback")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Feedback created successfully",
		"data":    data,
	})

}
