package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"errors"
	"gorm.io/gorm"
	"regexp"
)

func CreateUserService(user *models.User) (*models.User, error) {
	var existingUser models.User

	// Check if user with email already exists
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := utils.GetHashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func LoginUserService(loginData models.LoginModel) (*models.User, error) {
	var user models.User

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	query := database.DB
	if emailRegex.MatchString(loginData.Identifier) {
		query = query.Where("email = ?", loginData.Identifier)
	} else {
		query = query.Where("name = ?", loginData.Identifier)
	}

	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	valid := utils.ValidatePassword(loginData.Password, user.Password)
	if !valid {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func GetUserService(id string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
