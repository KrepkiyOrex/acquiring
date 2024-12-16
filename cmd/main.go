package main

import (
	"log"
	"sync"

	"github.com/KrepkiyOrex/acquiring/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/database/postgres"
	"github.com/KrepkiyOrex/acquiring/internal/service"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("========== Docker acquiring started ... ===========")

	var wg sync.WaitGroup
	wg.Add(2)

	var db1, db2 *gorm.DB
	var err1, err2 error

	go func() {
		defer wg.Done()
		db1, err1 = postgres.ConnectToDB("acquiring")
		if err1 != nil {
			log.Fatalf("Failed to connect to the database acquiring: %v", err1)
		}
	}()

	go func() {
		defer wg.Done()
		// db2, err := postgres.ConnectBank()
		db2, err2 = postgres.ConnectToDB("bank")
		if err2 != nil {
			log.Fatalf("Failed to connect to the database bank: %v", err2)
		}
	}()

	wg.Wait()

	if err := db1.AutoMigrate(&service.Transactions{}); err != nil {
		log.Fatalf("Could not migrate the database: %v", err)
	}

	db1 = db1.Debug()

	app := fiber.New()

	transactionService := service.NewAcquiringRepos(db1) // DB *gorm.DB

	bankRepos := service.NewBankRepos(db2) // DB *gorm.DB
	bankService := service.NewService(bankRepos)

	handlers.SetupRoutes(app, transactionService, bankService)

	// producer.Producer()

	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
