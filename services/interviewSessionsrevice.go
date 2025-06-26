package services

import (
	"backend/database"
	"backend/models"
	"errors"
	"time"
)

// StartInterviewSession starts a new interview session for a specific user and domain.
func StartInterviewSession(interviewSession *models.InterviewSession, user_id string) (*models.InterviewSession, error) {
	// Check if the user exists
	var user models.User
	if err := database.DB.First(&user, "id = ?", user_id).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Check if the domain exists
	var domain models.Domain
	if err := database.DB.First(&domain, "id = ?", interviewSession.DomainID).Error; err != nil {
		return nil, errors.New("domain not found")
	}

	// Set associations
	interviewSession.UserID = user.ID
	interviewSession.DomainID = domain.ID
	interviewSession.StartedAt = time.Now()

	// Create session
	if err := database.DB.Create(&interviewSession).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Preload("Domain").
		First(&interviewSession, "id = ?", interviewSession.ID).Error; err != nil {
		return nil, err
	}

	return interviewSession, nil
}

func EndInterviewSession(sessionID string) (*models.InterviewSession, error) {
	var interviewSession models.InterviewSession

	// Find the interview session by ID
	if err := database.DB.First(&interviewSession, "id = ?", sessionID).Error; err != nil {
		return nil, errors.New("interview session not found")
	}

	// Save the changes to the database
	interviewSession.CompletedAt = time.Now()
	if err := database.DB.Save(&interviewSession).Error; err != nil {
		return nil, err
	}

	return &interviewSession, nil
}
func GetAllInterviewSession(user_id string) ([]models.InterviewSession, error) {
	var interviewSessions []models.InterviewSession
	if err := database.DB.Preload("Domain").Where("user_id = ?", user_id).Find(&interviewSessions).Error; err != nil {
		return nil, err
	}
	return interviewSessions, nil
}

func GetInterviewSessionById(sessionID string) (*models.InterviewSession, error) {
	var interviewSession models.InterviewSession
	if err := database.DB.Preload("Domain").Preload("UserQuestions").First(&interviewSession, "id = ?", sessionID).Error; err != nil {
		return nil, errors.New("interview session not found")
	}
	return &interviewSession, nil
}
