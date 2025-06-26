package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or error loading it")
	}
}

func GetDBConfig() string {
	return os.Getenv("DATABASE_URL")
}

func GetPort() string {
	return os.Getenv("PORT")
}
func GetAdminpassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
func GetBAckedUrl() string {
	return os.Getenv("BACKEND_URL")
}

func GetJWTSecret(role string) string {
	if role == "admin" {
		return os.Getenv("ADMIN_JWT_SECRET")
	}
	if role == "user" {
		return os.Getenv("USER_JWT_SECRET")
	}

	return os.Getenv("JWT_SECRET")
}
