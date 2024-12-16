package handlers

import (
	"github.com/KrepkiyOrex/acquiring/internal/service"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, transService *service.TransRepos, svc *service.Service) {
// func SetupRoutes(app *fiber.App, transService *service.TransRepos, bankService *service.BankRepos) {
	api := app.Group("/api")

	api.Get("/get_allcard", svc.GetAllCardDetails)
	
	api.Get("/payment", service.ShowPaymentPage) // payment page
	api.Post("/deduct_balance", svc.DeductFromAccount)
	api.Post("/add_funds", svc.AddFunds)

	// api.Get("/pay", service.ShowCreatePage) // placebo
	api.Post("/create_transaction", transService.CreateTransaction)
	api.Delete("/delete_transaction/:id", transService.DeleteTransaction)
	api.Get("/get_transactions", transService.GetTransactions)
	api.Get("/get_transaction/:id", transService.GetTransByID)
}
