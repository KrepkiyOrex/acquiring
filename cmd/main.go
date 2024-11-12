package main

import (
	"acquiring/internal/postgres"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Структура DB
type DB struct {
	*postgres.DB
}

// Интерфейс Crudable
type Crudable interface {
	GetTransID() int64
	CreateTable(table string) string
	GetInsertQuery() (string, []interface{})
	GetUpdateQuery() string
	GetDeleteQuery() string
}

// Структура transactions
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

// Реализация метода GetTransID для структуры User
func (t *transactions) GetTransID() int64 {
	return t.TransactionID
}

// Реализация метода GetInsertQuery для структуры User
func (t *transactions) GetInsertQuery() (string, []interface{}) {
	// createdAt := time.Now().Truncate(time.Second) // Убираем микросекунды
	query := `INSERT INTO transactions
				(order_id, user_id, amount, currency, status, payment_method)
				VALUES ($1, $2, $3, $4, $5, $6)`
	return query, []interface{}{t.OrderID, t.UserID, t.Amount, t.Currency, t.Status, t.PaymentMethod}
}

// Реализация метода GetUpdateQuery для структуры User
func (t *transactions) GetUpdateQuery() string {
	return fmt.Sprintf("UPDATE users SET name = '%s', email = '%s' WHERE id = %d")
}

// Реализация метода GetDeleteQuery для структуры User
func (t *transactions) GetDeleteQuery() string {
	return fmt.Sprintf("DELETE FROM users WHERE id = %d")
}

// Функция для добавления записи
func Add[T Crudable](db *DB, entity T) error {
	query, args := entity.GetInsertQuery()
	_, err := db.Exec(query, args...)
	return err
}

// Функция для обновления записи
func Update[T Crudable](db *DB, entity T) error {
	query := entity.GetUpdateQuery()
	_, err := db.Exec(query)
	return err
}

// Функция для удаления записи
func Delete[T Crudable](db *DB, entity T) error {
	query := entity.GetDeleteQuery()
	_, err := db.Exec(query)
	return err
}

func Create[T Crudable](db *DB, entity T, table string) error {
	query := entity.CreateTable(table)
	_, err := db.Exec(query)
	return err
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
	// CurrencyForm := r.FormValue("currency")
	// StatusForm := r.FormValue("status")
	// PaymentMethodForm := r.FormValue("payment_method")

	t.OrderID = 465
	t.UserID = 98
	t.Amount = 2356.54
	// t.Currency = CurrencyForm
	// t.Status = StatusForm
	// t.PaymentMethod = PaymentMethodForm
}

func transactionHandler(w http.ResponseWriter, r *http.Request) {
	transac := &transactions{}

	transac.SetData(r)

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	defer db.Close()

	dbInstance := &DB{db}

	err = Add(dbInstance, transac)
	if err != nil {
		log.Fatal("Ошибка при добавлении транзакции:", err)
	}

	fmt.Println("Транзакция успешно добавлена!")

}

func wel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hey, hey everybody!")
	fmt.Fprintln(w, "Peu, peu, peu!")
}

func main() {
	http.HandleFunc("/", transactionHandler)

	// http.HandleFunc("/", wel)

	log.Println("Server started...")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
