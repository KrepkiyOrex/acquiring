package main

import (
	"log"

	"github.com/KrepkiyOrex/acquiring/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/database/postgres"
	"github.com/KrepkiyOrex/acquiring/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("===== Docker acquiring started ... =====")

	db, err := postgres.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(&service.Transactions{}); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	db = db.Debug()

	transactionService := service.NewRepository(db) // DB *gorm.DB
	app := fiber.New()

	handlers.SetupRoutes(app, transactionService)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
