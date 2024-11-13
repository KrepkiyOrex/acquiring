package service

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/KrepkiyOrex/acquiring/internal/postgres"

	log "github.com/sirupsen/logrus"
)

var transac = &transactions{}

type transactions struct {
	TransactionID int64
	OrderID       int64
	UserID        int64
	Amount        float64
	Currency      string
	Status        string
	PaymentMethod string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Интерфейс Crudable
type Crudable interface {
	GetTransID() int64
	CreateTable(table string) string
	GetInsertQuery() (string, []interface{})
	GetUpdateQuery() string
	GetDeleteQuery() string
}

func (t *transactions) GetTransID() int64 {
	return t.TransactionID
}

func (t *transactions) GetInsertQuery() (string, []interface{}) {
	query := `INSERT INTO transactions
				(order_id, user_id, amount, currency, status, payment_method)
				VALUES ($1, $2, $3, $4, $5, $6)`
	return query, []interface{}{t.OrderID, t.UserID, t.Amount, t.Currency, t.Status, t.PaymentMethod}
}

func (t *transactions) GetUpdateQuery() string {
	return fmt.Sprintf("UPDATE users SET name = '%s', email = '%s' WHERE id = %d")
}

func (t *transactions) GetDeleteQuery() string {
	return fmt.Sprintf("DELETE FROM users WHERE id = %d")
}

// func (t *transactions) Insert(db *DB) {
func (t *transactions) Insert() {
	query, args := t.GetInsertQuery()
	_, err := postgres.GetDB().Exec(query, args...)
	if err != nil {
		log.Fatal("Ошибка при добавлении данных:", err)
	}
}

func (t *transactions) CreateTable(table string) string {
	return fmt.Sprintf("CREATE TABLE %v (id INT, name VARCHAR(255), age INT)", table)
}

func (t *transactions) parseOrderID(r *http.Request) {
	orderIDStr := r.FormValue("order_id")
	var err error

	t.OrderID, err = strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		log.Fatal("Error while parsing OrderID:", err)
	}
}

func (t *transactions) parseUserID(r *http.Request) {
	userIDStr := r.FormValue("user_id")
	var err error

	t.UserID, err = strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Fatal("Error while parsing UserID:", err)
	}
}

func (t *transactions) parseAmont(r *http.Request) {
	AmountStr := r.FormValue("amount")
	var err error

	t.Amount, err = strconv.ParseFloat(AmountStr, 64)
	if err != nil {
		log.Fatal("Error while parsing Amount:", err)
	}
}

func (t *transactions) SetData(r *http.Request) {
	CurrencyForm := r.FormValue("currency")
	StatusForm := r.FormValue("status")
	PaymentMethodForm := r.FormValue("payment_method")

	t.parseOrderID(r)
	t.parseUserID(r)
	t.parseAmont(r)

	t.Currency = CurrencyForm
	t.Status = StatusForm
	t.PaymentMethod = PaymentMethodForm
}

func CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Error("Invalid request method: ", http.StatusMethodNotAllowed)
	}

	transac.SetData(r)

	transac.Insert()

	log.Info("Transaction created:", transac)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Transaction created successfully")
}

func PaymentPage(w http.ResponseWriter, r *http.Request) {
	page := []string{"main.html"}

	template, err := template.ParseFiles(page...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Wel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hey, hey everybody!")
	fmt.Fprintln(w, "Peu, peu, peu!")

	fmt.Fprintln(w, "/start track synthetik legion rising cyberpunk combat 2")
}

// func TransactionHandler(w http.ResponseWriter, r *http.Request) {
// 	transac := &transactions{}

// 	transac.SetData(r)

// 	db, err := postgres.Connect()
// 	if err != nil {
// 		log.Fatal("Не удалось подключиться к базе данных:", err)
// 	}
// 	defer db.Close()

// 	dbInstance := &DB{db}

// 	err = Add(dbInstance, transac)
// 	if err != nil {
// 		log.Fatal("Ошибка при добавлении транзакции:", err)
// 	}

// 	fmt.Println("Транзакция успешно добавлена!")
// }
