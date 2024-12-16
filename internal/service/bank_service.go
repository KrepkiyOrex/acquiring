package service

import (
	"net/http"

	"github.com/KrepkiyOrex/acquiring/internal/producer"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardData struct {
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

// func (AccountBalance) TableName() string {
// 	return "card_data"
// }

func NewBankRepos(db *gorm.DB) *BankRepos {
	return &BankRepos{DB: db}
}

func ShowPaymentPage(c *fiber.Ctx) error {
	return c.Render("internal/source/payment.html", fiber.Map{
		"ErrorMessage": "",
	})
}

// вычитаем из баланса
func (card CardData) DecrementBalance() clause.Expr {
	return gorm.Expr("balance - ?", card.Balance)
}

// пополняем баланс
func (card CardData) IncrementBalance() clause.Expr {
	return gorm.Expr("balance + ?", card.Balance)
}

func (bank *BankRepos) DeductFromAccount(ctx *fiber.Ctx) error {
	details := &CardData{}

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	result := bank.DB.Model(&CardData{}).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Where("expiry_date = ?", details.ExpiryDate).
		Where("encrypted_CVV = ?", details.EncryptedCvv).
		Where("balance >= ?", details.Balance). // проверка, хватает ли денег
		Update("balance", details.DecrementBalance())

	transaction := Transactions{
		TransactionID: 123456,
		OrderID:       78910,  // получаем с магазина (другого приложения)
		UserID:        101112, // получаем с магазина (другого приложения)
		Amount:        1000.50,
	}

	transaction.SetBalance(details.Balance)

	producer.Producer()

	// создай метод, внутри которого будут методы, что установят все эти
	// значения с других методово или программ

	// если RowsAffected равно 0, это значит, что либо карта не найдена,
	// либо на счете недостаточно средств.
	if result.RowsAffected == 0 {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Not enough money or card not found",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Balance has been successfull deducted"})
}

func (tsn *Transactions) SetBalance(balance float64) *Transactions {
	tsn.Amount = balance
	return tsn
}

func (bank *BankRepos) AddFunds(ctx *fiber.Ctx) error {
	details := &CardData{}

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	result := bank.DB.Model(&CardData{}).
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

func (bank *BankRepos) GetAllCardDetails(ctx *fiber.Ctx) error {
	balanceModels := &[]CardData{}
	if err := bank.DB.Find(&balanceModels).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get cards"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "cards fetched successfully",
		"data":    balanceModels,
	})
	return nil
}
