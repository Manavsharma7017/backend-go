package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"errors"
	"gorm.io/gorm"
	"regexp"
)

func CreateAdminUserService(user *models.AdminUser) (*models.AdminUser, error) {

	var existingUser models.AdminUser
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("admin user with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("admin user with this username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

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

func LoginAdminUserService(adminuserlogin models.LoginModel) (*models.AdminUser, error) {
	var user models.AdminUser
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	query := database.DB
	if emailRegex.MatchString(adminuserlogin.Identifier) {
		query = query.Where("email = ?", adminuserlogin.Identifier)
	} else {
		query = query.Where("username = ?", adminuserlogin.Identifier)
	}
	if err := query.First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	err := utils.ValidatePassword(adminuserlogin.Password, user.Password)
	if !err {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}

func GetAllAdminUserService() ([]models.AdminUser, error) {
	var users []models.AdminUser
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func GetAdminDashboardService() (map[string]int64, error) {
	var userCount int64
	if err := database.DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return nil, err
	}

	var domainCount int64
	if err := database.DB.Model(&models.Domain{}).Count(&domainCount).Error; err != nil {
		return nil, err
	}

	var questionCount int64
	if err := database.DB.Model(&models.Question{}).Count(&questionCount).Error; err != nil {
		return nil, err
	}
	var sessioncount int64
	if err := database.DB.Model(&models.InterviewSession{}).Count(&sessioncount).Error; err != nil {
		return nil, err
	}

	dashboardData := map[string]int64{
		"total_users":     userCount,
		"total_domains":   domainCount,
		"total_questions": questionCount,
		"total_sessions":  sessioncount,
	}

	return dashboardData, nil
}
