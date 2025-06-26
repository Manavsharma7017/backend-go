package services

import (
	"backend/database"
	"backend/models"
	"errors"
)

func CreateUserQuestion(userQuestion *models.UserQuestion, user_id string) (*models.UserQuestion, error) {
	// Set user ID if not already set
	if userQuestion.UserID == "" {
		userQuestion.UserID = user_id
	}

	// Create the user question entry in the database
	if err := database.DB.Create(&userQuestion).Error; err != nil {
		return nil, errors.New("failed to create user question")
	}

	// Fetch the record with all associations preloaded
	if err := database.DB.
		Preload("User").       // Load the associated User
		Preload("Session").    // Load the associated Session
		Preload("UserDomain"). // Load the associated UserDomain
		First(&userQuestion, "id = ?", userQuestion.ID).Error; err != nil {
		return nil, errors.New("failed to fetch created user question with associations")
	}

	return userQuestion, nil
}

func GetAllUserQuestions() ([]models.UserQuestion, error) {
	var userQuestions []models.UserQuestion
	if err := database.DB.Preload("User").Preload("Session").Preload("UserDomain").Find(&userQuestions).Error; err != nil {
		return nil, errors.New("failed to retrieve user questions")
	}
	return userQuestions, nil
}

func GetUserQuestionByID(id string) (*models.UserQuestion, error) {
	var userQuestion models.UserQuestion
	if err := database.DB.Preload("Responses").First(&userQuestion, "id = ?", id).Error; err != nil {
		return nil, errors.New("failed to retrieve user question by ID")
	}
	return &userQuestion, nil
}
func DeleteUserQuestion(id string) error {
	var userQuestion models.UserQuestion
	if err := database.DB.First(&userQuestion, "id = ?", id).Error; err != nil {
		return errors.New("user question not found")
	}

	if err := database.DB.Delete(&userQuestion).Error; err != nil {
		return errors.New("failed to delete user question")
	}

	return nil
}
