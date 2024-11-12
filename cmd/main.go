package main

import (
	"log"

	"github.com/KrepkiyOrex/acquiring/internal/payment/handlers"
)

func main() {

	log.Println("===== Payment server welcomes you! =====")

	handlers.StartServer()

}
