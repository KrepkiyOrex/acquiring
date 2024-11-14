package main

import (
	"log"

	"github.com/KrepkiyOrex/acquiring/internal/payment/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/payment/service"
)

func main() {

	log.Println("===== Payment server welcomes you! =====")

	// postgres.Connect()
	_, err := service.Connect()
	if err != nil {
		log.Fatal("Error connect: ", err)
	}

	defer service.GetDB().Close()

	handlers.StartServer()
}
