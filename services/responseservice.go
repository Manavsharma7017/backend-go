package services

import (
	"backend/database"
	"backend/models"
)

func CreateResponse(response *models.Response) (*models.Response, error) {
	// Step 1: Create the response in the database
	if err := database.DB.Create(response).Error; err != nil {
		return nil, err
	}

	// Step 2: Preload related associations (if needed)
	if err := database.DB.
		Preload("Session").
		Preload("UserQuestion").
		First(response, "responseid = ?", response.ResponseID).Error; err != nil {
		return nil, err
	}

	return response, nil
}

func GetResponse(responseID string) (*models.Response, error) {
	if responseID == "" {
		return nil, nil
	}
	var response models.Response
	if err := database.DB.Preload("Feedback").Where("responseid = ?", responseID).First(&response).Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func GetSessionResponse(sessionID string) ([]models.Response, error) {
	if sessionID == "" {
		return nil, nil
	}
	var responses []models.Response
	if err := database.DB.Where("session_id = ?", sessionID).Find(&responses).Error; err != nil {
		return nil, err
	}
	return responses, nil
}
