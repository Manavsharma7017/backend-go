package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"sync"

	"github.com/gofiber/fiber/v2"
)

func StartInterviewSession(c *fiber.Ctx) error {
	var insterviewsession models.InterviewSession
	var user_id = c.Locals("user_id")
	if user_id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	if err := c.BodyParser(&insterviewsession); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := utils.ValidateStruct(insterviewsession); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload validrtion",
		})
	}
	session, err := services.StartInterviewSession(&insterviewsession, user_id.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start interview session",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Interview session started successfully",
		"session": fiber.Map{
			"id":          session.ID,
			"domain_id":   session.DomainID,
			"domain_name": session.Domain.Name,
			"description": session.Description,
			"user_id":     session.UserID,
			"started_at":  session.StartedAt,
		},
	})

}

func EndInterviewSession(c *fiber.Ctx) error {
	// Assuming you have a way to identify the interview session to end, e.g., by ID
	sessionId := c.Params("id")
	if sessionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Logic to end the interview session goes here
	session, err := services.EndInterviewSession(sessionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to end interview session",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Interview session ended successfully",
		"session": fiber.Map{
			"id":        session.ID,
			"domain_id": session.DomainID,
			"user_id":   session.UserID,
			"ended_at":  session.CompletedAt,
		},
	})

}
func GetAllInterviewSession(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")

	if user_id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	interviewSessions, err := services.GetAllInterviewSession(user_id.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get interview sessions",
		})
	}
	interviewSessionChan := make(chan models.InterviewSession)
	var sanitizedata []fiber.Map
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for interviewSession := range interviewSessionChan {
				sanitized := fiber.Map{
					"id":           interviewSession.ID,
					"domain_id":    interviewSession.DomainID,
					"domain_name":  interviewSession.Domain.Name,
					"description":  interviewSession.Description,
					"started_at":   interviewSession.StartedAt,
					"completed_at": interviewSession.CompletedAt,
				}
				mu.Lock()
				sanitizedata = append(sanitizedata, sanitized)
				mu.Unlock()
			}
		}()
	}
	for _, interviewSession := range interviewSessions {
		interviewSessionChan <- interviewSession
	}
	close(interviewSessionChan)
	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"interview_sessions": sanitizedata,
	})

}
func GetInterviewSessionById(c *fiber.Ctx) error {
	sessionId := c.Params("id")
	if sessionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}
	insterviewsession, err := services.GetInterviewSessionById(sessionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get interview session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Interview started",
		"session": fiber.Map{
			"id":          insterviewsession.ID,
			"domain_id":   insterviewsession.DomainID,
			"domain_name": insterviewsession.Domain.Name,
			"description": insterviewsession.Description,
			"Questions":   insterviewsession.UserQuestions,
		},
	})

}
