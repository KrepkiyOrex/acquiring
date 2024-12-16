package producer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
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

func connRetryToKafka(addr string) *kafka.Writer {
	for {
		writer := &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Topic:    "transaction",
			Balancer: &kafka.LeastBytes{},
		}

		// Проверка подключения через создание объекта Writer
		log.Println("Trying to connect to Kafka...")
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// checking connect
		err := writer.WriteMessages(ctx, kafka.Message{})
		if err == nil {
			log.Println("Successfully connected to Kafka.")
			return writer
		}

		log.Printf("Failed to connect to Kafka: %v. Retrying in 3 seconds...\n", err)
		writer.Close()
		time.Sleep(3 * time.Second)
	}
}

func Producer() {
	writer := connRetryToKafka("kafka:9093")
	defer writer.Close()

	for {
		transaction := Transactions{
			TransactionID: 123456,
			OrderID:       78910, // получаем с магазина (другого приложения)
			UserID:        101112, // получаем с магазина (другого приложения)
			Amount:        1000.50, 
			Currency:      "USD", // - по умолчанию будет статичен
			Status:        "SUCCESS", // тут будет генерировать в DeductFromAccount
			PaymentMethod: "CREDIT_CARD", // получаем с магазина (другого приложения)
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		data, err := json.Marshal(transaction)
		if err != nil {
			log.Fatal("could not serialize transaction: ", err)
		}

		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(strconv.FormatInt(transaction.TransactionID, 10)),
				Value: data,
			},
		)

		if err != nil {
			log.Fatal("could not write message: ", err)
		}

		log.Printf("Sent transaction: %s", data)

		time.Sleep(3 * time.Second) // задержка для следующей итерации
	}
}
