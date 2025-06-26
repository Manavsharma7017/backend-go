package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateQuestionHandler(c *fiber.Ctx) error {
	var Question models.Question
	if err := c.BodyParser(&Question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}
	if err := utils.ValidateStruct(Question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})

	}
	_, err := services.CreateQuestion(&Question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Question created successfully",
		"question": fiber.Map{
			"id":         Question.ID,
			"question":   Question.Text,
			"domain_id":  Question.DomainID,
			"created_at": Question.CreatedAt,
			"updated_at": Question.UpdatedAt,
		},
	})

}
func GetAllQuestionsHandler(c *fiber.Ctx) error {
	questions, err := services.GetAllQuestions()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
	})
}
func GetQuestionFilterHandler(c *fiber.Ctx) error {
	difficulty := c.Query("difficulty")
	domainID := c.QueryInt("domain_id")

	questions, err := services.GetQuestions(difficulty, domainID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"questions": questions,
	})
}

func GetQuestionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	question, err := services.GetQuestionByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"question": question,
	})
}

func UpdateQuestionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var question models.Question

	if err := c.BodyParser(&question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}

	if err := utils.ValidateStruct(question); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	updatedQuestion, err := services.UpdateQuestion(id, &question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Question updated successfully",
		"question": updatedQuestion,
	})
}

func DeleteQuestionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	err := services.DeleteQuestion(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Question deleted successfully",
	})
}
