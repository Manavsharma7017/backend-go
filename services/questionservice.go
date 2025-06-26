package services

import (
	"backend/database"
	"backend/models"
	"errors"

	"gorm.io/gorm"
)

func GetAllQuestions() ([]models.Question, error) {
	var questions []models.Question
	if err := database.DB.Preload("Domain").Find(&questions).Error; err != nil {
		return nil, errors.New("failed to fetch questions")
	}
	return questions, nil
}

func CreateQuestion(question *models.Question) (*models.Question, error) {
	var domain models.Domain

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Step 1: Check if Domain exists
		if err := tx.First(&domain, "id = ?", question.DomainID).Error; err != nil {
			return errors.New("domain not found")
		}

		// Step 2: Create Question
		if err := tx.Create(question).Error; err != nil {
			return errors.New("failed to create question")
		}

		// Step 3: Reload Question with Domain preloaded
		if err := tx.Preload("Domain").First(question, "id = ?", question.ID).Error; err != nil {
			return errors.New("failed to retrieve created question with domain")
		}

		return nil // if everything succeeds, commit
	})

	if err != nil {
		return nil, err // rollback will happen automatically if err != nil
	}

	return question, nil
}

func GetQuestions(difficulty string, domainID int) ([]models.Question, error) {
	var questions []models.Question
	query := database.DB.Preload("Domain")

	if difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}

	if domainID != 0 {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Find(&questions).Error; err != nil {
		return nil, errors.New("failed to fetch questions")
	}

	return questions, nil
}

func GetQuestionByID(id string) (*models.Question, error) {
	var question models.Question
	if err := database.DB.Preload("Domain").First(&question, "id = ?", id).Error; err != nil {
		return nil, errors.New("question not found")
	}

	return &question, nil
}

func UpdateQuestion(id string, updatedData *models.Question) (*models.Question, error) {
	var existingQuestion models.Question
	if err := database.DB.First(&existingQuestion, "id = ?", id).Error; err != nil {
		return nil, errors.New("question not found")
	}

	// Update fields
	existingQuestion.Text = updatedData.Text
	existingQuestion.Difficulty = updatedData.Difficulty
	existingQuestion.DomainID = updatedData.DomainID

	if err := database.DB.Save(&existingQuestion).Error; err != nil {
		return nil, errors.New("failed to update question")
	}

	if err := database.DB.Preload("Domain").First(&existingQuestion, "id = ?", id).Error; err != nil {
		return nil, errors.New("failed to retrieve updated question")
	}

	return &existingQuestion, nil
}

func DeleteQuestion(id string) error {
	if err := database.DB.Delete(&models.Question{}, "id = ?", id).Error; err != nil {
		return errors.New("failed to delete question")
	}
	return nil
}
