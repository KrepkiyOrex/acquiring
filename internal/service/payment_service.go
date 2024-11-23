package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Transactions struct {
	TransactionID int64     `json:"transaction_id" gorm:"primaryKey"`
	OrderID       int64     `json:"orderId"`
	UserID        int64     `json:"userId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"paymentMethod"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransRepos struct {
	DB *gorm.DB
}

func ShowCreatePage(c *fiber.Ctx) error {
	return c.Render("internal/source/form.html", fiber.Map{
		"ErrorMessage": "",
	})
}

func NewAcquiringRepos(db *gorm.DB) *TransRepos {
	return &TransRepos{DB: db}
}

func (r *TransRepos) CreateTransaction(ctx *fiber.Ctx) error {
	transaction := Transactions{}

	err := ctx.BodyParser(&transaction)
	if err != nil {
		return ctx.Render("internal/source/form.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	// if ok, err := ValidateTrans(transaction); !ok {
	// 	return ctx.Render("internal/source/form.html", fiber.Map{
	// 		"ErrorMessage": err.Error(),
	// 	})
	// }

	if err := r.DB.Create(&transaction).Error; err != nil {
		return ctx.Render("internal/source/form.html", fiber.Map{
			"ErrorMessage": "Could not create transaction",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Transaction has been added"})
}

func (r *TransRepos) GetTransByID(context *fiber.Ctx) error {

	id := context.Params("id")
	transModel := &Transactions{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("transaction_id = ?", id).First(transModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the transactions"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "transactions id fetched successfully",
		"data":    transModel,
	})
	return nil
}

func (r *TransRepos) GetTransactions(c *fiber.Ctx) error {
	transModels := &[]Transactions{}
	if err := r.DB.Find(&transModels).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get transactions"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "transactions fetched successfully",
		"data":    transModels,
	})
	return nil
}

func (r *TransRepos) DeleteTransaction(c *fiber.Ctx) error {
	transModel := Transactions{}
	txnID := c.Params("id")
	if txnID == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "transaction_id cannot be empty"})
		return nil
	}

	if err := r.DB.Delete(&transModel, txnID).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete transaction"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "transaction deleted successfully"})
	return nil
}
