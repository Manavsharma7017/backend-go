package services

import (
	"backend/database"
	"backend/models"
	"errors"

	"gorm.io/gorm"
)

func CreateUserDomain(userdomain *models.UserDomain) (*models.UserDomain, error) {
	var existingUserDomain models.UserDomain
	err := database.DB.
		Where("user_id = ? AND domain_id = ?", userdomain.UserID, userdomain.DomainID).
		First(&existingUserDomain).Error

	if err == nil {
		// ✅ Found existing one
		return &existingUserDomain, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// ❌ Some other DB error
		return nil, err
	}

	// ✅ Not found → create new
	if err := database.DB.Create(userdomain).Error; err != nil {
		return nil, err
	}
	return userdomain, nil
}

func GetAllUserDomains(userdomains *[]models.UserDomain, user_id string) (*[]models.UserDomain, error) {
	if user_id == "" {
		return nil, errors.New("user_id is unauthorized")
	}
	if err := database.DB.Preload("Domain").Where("user_id = ?", user_id).Find(userdomains).Error; err != nil {
		return nil, err
	}

	return userdomains, nil
}
func GetUserDomainByID(userdomain *models.UserDomain, id string) (*models.UserDomain, error) {
	if id == "" {
		return nil, errors.New("user domain ID is required")
	}
	if err := database.DB.Where("id = ?", id).First(userdomain).Error; err != nil {
		return nil, err
	}
	return userdomain, nil
}
func DeletUserDomainByID(userdomain *models.UserDomain, id string) (*models.UserDomain, error) {
	if id == "" {
		return nil, errors.New("user domain ID is required")
	}
	if err := database.DB.Where("id = ?", id).Delete(userdomain).Error; err != nil {
		return nil, err
	}
	if err := database.DB.Where("id = ?", id).First(userdomain).Error; err == nil {
		return nil, errors.New("user domain not deleted")
	}
	return userdomain, nil
}
