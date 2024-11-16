package handlers

import (
	"github.com/KrepkiyOrex/acquiring/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, transactionService *service.Repository) {
	api := app.Group("/api")
	api.Post("/create_transaction", transactionService.CreateTransaction)
	api.Delete("/delete_transaction/:id", transactionService.DeleteTransaction)
	api.Get("/get_transaction", transactionService.GetTransactions)
}
