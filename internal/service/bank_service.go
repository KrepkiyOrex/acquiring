package service

import (
	"fmt"
	"net/http"

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

type BankRepository interface {
	DeductFromAccount(ctx *fiber.Ctx, details *CardData) error
	AddFunds(ctx *fiber.Ctx) error
	GetAllCardDetails(ctx *fiber.Ctx) error
}

// type Application struct {
// 	db dbContract
// }

func NewBankRepos(db *gorm.DB) *BankRepos {
	return &BankRepos{DB: db}
}

func ShowPaymentPage(ctx *fiber.Ctx) error {
	return ctx.Render("internal/source/payment.html", fiber.Map{
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

func (bank *BankRepos) DeductFromAccount(ctx *fiber.Ctx, details *CardData) error {
	// details := &CardData{}

	// if err := ctx.BodyParser(details); err != nil {
	// 	return ctx.Render("internal/source/payment.html", fiber.Map{
	// 		"ErrorMessage": "Failed to parse request",
	// 	})
	// }
	fmt.Println("Deduct: ", details)

	result := bank.DB.Model(&CardData{}).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Where("expiry_date = ?", details.ExpiryDate).
		Where("encrypted_CVV = ?", details.EncryptedCvv).
		Where("balance >= ?", details.Balance). // проверка, хватает ли денег
		Update("balance", details.DecrementBalance())

	// если RowsAffected равно 0, это значит, что либо карта не найдена,
	// либо на счете недостаточно средств.
	if result.RowsAffected == 0 {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Not enough money or card not found",
		})
	}

	// producer.Producer()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Balance has been successfull deducted"})
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
			"ErrorMessage": "Not enough money or card not found"})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Funds added successfully"})
}

func (tr *Transactions) SetAmount(details CardData) {
	tr.Amount = details.Balance
}

func NewTransaction() *Transactions {
	return &Transactions{}
}

func ProcessPayment(bankRepo BankRepository, transRepo *TransRepos, ctx *fiber.Ctx) error {
	details := &CardData{}

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("internal/source/payment.html", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	if err := bankRepo.DeductFromAccount(ctx, details); err != nil {
		return err
	}

	transaction := NewTransaction()

	transaction.SetAmount(*details)

	if err := transRepo.CreateTransaction(ctx, transaction); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Payment processed successfully",
	})
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

type Service struct {
	BankRepo BankRepository
}

func NewService(repo BankRepository) *Service {
	return &Service{BankRepo: repo}
}

func (s *Service) AddFunds(ctx *fiber.Ctx) error {
	return s.BankRepo.AddFunds(ctx)
}

func (s *Service) DeductFromAccount(ctx *fiber.Ctx, details *CardData) error {
	return s.BankRepo.DeductFromAccount(ctx, details)
}

func (s *Service) GetAllCardDetails(ctx *fiber.Ctx) error {
	return s.BankRepo.GetAllCardDetails(ctx)
}
