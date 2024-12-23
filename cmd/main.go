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

	db1, db2 := postgres.SetupDataBases()

	defer db1.Close()
	defer db2.Close()

	if err := db1.AutoMigrate(&service.Transactions{}); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	db1 = db1.Debug()

	app := fiber.New()

	transactionService := service.NewTransRepos(db1.DB) // DB *gorm.DB

	bankRepos := service.NewBankRepos(db2.DB) // DB *gorm.DB
	service := service.NewService(bankRepos, transactionService)

	handlers.SetupRoutes(app, service)

	// producer.Producer()

	if err := app.Listen(":8081"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
