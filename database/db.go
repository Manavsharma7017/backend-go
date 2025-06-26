package database

import (
	"backend/config"
	"backend/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.GetDBConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to database")
	DB = db
}
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.AdminUser{},
		&models.Domain{},
		&models.Question{},
		&models.InterviewSession{},
		&models.Response{},
		&models.Feedback{},
		&models.UserDomain{},
		&models.UserQuestion{},
	)
	if err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	fmt.Println("✅ Database migrated")
}
