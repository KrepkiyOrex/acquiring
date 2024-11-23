package main

import (
	"log"

	"github.com/KrepkiyOrex/acquiring/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/database/postgres"
	"github.com/KrepkiyOrex/acquiring/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("========== Docker acquiring started ... ===========")

	db1, err := postgres.ConnectTrans()
	if err != nil {
		log.Fatalf("Failed to connect to the database acquiring: %v", err)
	}

	db2, err := postgres.ConnectBank()
	if err != nil {
		log.Fatalf("Failed to connect to the database bank: %v", err)
	}

	if err := db1.AutoMigrate(&service.Transactions{}); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	db1 = db1.Debug()

	app := fiber.New()

	transactionService := service.NewAcquiringRepos(db1) // DB *gorm.DB

	bankService := service.NewBankRepos(db2) // DB *gorm.DB

	handlers.SetupRoutes(app, transactionService, bankService)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
