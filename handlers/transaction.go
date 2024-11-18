package handlers

import (
	"github.com/KrepkiyOrex/acquiring/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, transService *service.Repository) {
	api := app.Group("/api")
	api.Get("/pay", service.ShowCreatePage) // payment page
	api.Post("/create_transaction", transService.CreateTransaction)

	api.Delete("/delete_transaction/:id", transService.DeleteTransaction)
	api.Get("/get_transactions", transService.GetTransactions)
	api.Get("/get_transaction/:id", transService.GetTransByID)
}
