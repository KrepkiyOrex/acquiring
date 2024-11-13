package main

import (
	"log"

	"github.com/KrepkiyOrex/acquiring/internal/payment/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/postgres"
)

func main() {

	log.Println("===== Payment server welcomes you! =====")

	postgres.Connect()
	defer postgres.GetDB().Close()

	handlers.StartServer()
}
