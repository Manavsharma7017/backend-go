package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, Admin-Password",
	}))
	database.ConnectDB()
	database.AutoMigrate()
	PDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get DB instance: %v", err)
	}
	defer func() {
		if err := PDB.Close(); err != nil {
			log.Fatalf("❌ Failed to close DB: %v", err)
		}
	}()
	routes.AdminUserRoutes(app)
	routes.DomainRoutes(app)
	routes.FeedBackRoute(app)
	routes.InterviewSessionRoute(app)
	routes.QuestionRoutes(app)
	routes.ResponseRoute(app)
	routes.UserRoute(app)
	routes.UserQuestionRoute(app)
	routes.UserDomainRoute(app)

	port := config.GetPort()
	log.Fatal(app.Listen(port))

}
