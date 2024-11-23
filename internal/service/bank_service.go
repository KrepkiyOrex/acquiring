package service

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountBalance struct {
	ID                  int64   `json:"id" gorm:"primaryKey"`
	Balance             float64 `json:"balance"`
	EncryptedCardNumber string  `json:"encryptedCardNumber"`
	ExpiryDate          string  `json:"expiryDate"`
	EncryptedCvv        string  `json:"encryptedCvv"`
	EncryptedCardName   string  `json:"encryptedCardName"`
}

type BankRepos struct {
	DB *gorm.DB
}

func (AccountBalance) TableName() string {
	return "card_data"
}

func NewBankRepos(db *gorm.DB) *BankRepos {
	return &BankRepos{DB: db}
}

func ShowPaymentPage(c *fiber.Ctx) error {
	return c.Render("internal/source/payment.html", fiber.Map{
		"ErrorMessage": "",
	})
}

// вычитаем из баланса
func (a AccountBalance) DecrementBalance() clause.Expr {
	return gorm.Expr("balance - ?", a.Balance)
}

// пополняем баланс
func (a AccountBalance) IncrementBalance() clause.Expr {
	return gorm.Expr("balance + ?", a.Balance)
}

func (bank *BankRepos) DeductFromAccount(ctx *fiber.Ctx) error {
	details := &AccountBalance{}

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	result := bank.DB.Model(&AccountBalance{}).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Where("balance >= ?", details.Balance). // проверка, хватает ли денег
		Update("balance", details.DecrementBalance())

	// если RowsAffected равно 0, это значит, что либо карта не найдена,
	// либо на счете недостаточно средств.
	if result.RowsAffected == 0 {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Not enough money or card not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Balance has been successfull deducted"})
}

func (bank *BankRepos) AddFunds(ctx *fiber.Ctx) error {
	details := &AccountBalance{}

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	result := bank.DB.Model(&AccountBalance{}).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Update("balance", details.IncrementBalance())

	// если RowsAffected равно 0, это значит, что либо карта не найдена,
	// либо на счете недостаточно средств.
	if result.RowsAffected == 0 {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Not enough money or card not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Funds added successfully"})
}

func (r *BankRepos) GetAllCardDetails(c *fiber.Ctx) error {
	balanceModels := &[]AccountBalance{}
	if err := r.DB.Find(&balanceModels).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get cards"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "cards fetched successfully",
		"data":    balanceModels,
	})
	return nil
}
