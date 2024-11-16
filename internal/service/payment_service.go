package service

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Transactions struct {
	TransactionID int64     `json:"transaction_id" gorm:"primaryKey"`
	OrderID       int64     `json:"order_id"`
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (s *Repository) CreateTransaction(c *fiber.Ctx) error {
	transaction := Transactions{}

	if err := c.BodyParser(&transaction); err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	if err := s.DB.Create(&transaction).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create transaction"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "transaction has been added"})
	return nil
}

func (s *Repository) GetTransactions(c *fiber.Ctx) error {
	transModels := &[]Transactions{}
	if err := s.DB.Find(&transModels).Error; err != nil {
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

func (s *Repository) DeleteTransaction(c *fiber.Ctx) error {
	transModel := Transactions{}
	txnID := c.Params("id")
	if txnID == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "transaction_id cannot be empty"})
		return nil
	}

	if err := s.DB.Delete(&transModel, txnID).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete transaction"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "transaction deleted successfully"})
	return nil
}
