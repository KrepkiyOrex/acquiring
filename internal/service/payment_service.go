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

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func ShowCreatePage(c *fiber.Ctx) error {
	return c.Render("internal/source/form.html", fiber.Map{
		"ErrorMessage": "",
	})
}

func ShowPaymentPage(c *fiber.Ctx) error {
	return c.Render("internal/source/payment.html", fiber.Map{
		"ErrorMessage": "",
	})
}

func (r *Repository) CreateTransaction(ctx *fiber.Ctx) error {
	transaction := Transactions{}

	err := ctx.BodyParser(&transaction)
	if err != nil {
		return ctx.Render("internal/source/form.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	if ok, err := ValidateTrans(transaction); !ok {
		return ctx.Render("internal/source/form.html", fiber.Map{
			"ErrorMessage": err.Error(),
		})
	}

	

	if err := r.DB.Create(&transaction).Error; err != nil {
		return ctx.Render("internal/source/form.html", fiber.Map{
			"ErrorMessage": "Could not create transaction",
		})
	}

	// if r.DB.Update(&transaction).Error; err != nil {
	// 	return ctx.Render("internal/source/form.html", fiber.Map{
	// 		"ErrorMessage": "Could not create transaction",
	// 	})
	// }

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Transaction has been added"})
}

func ValidateTrans(trans Transactions) (bool, error) {
	if trans.Amount <= 1000 {
		return true, nil
	}
	return false, fmt.Errorf("not enough money")
}

func (r *Repository) GetTransByID(context *fiber.Ctx) error {

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

func (r *Repository) GetTransactions(c *fiber.Ctx) error {
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

func (r *Repository) DeleteTransaction(c *fiber.Ctx) error {
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
