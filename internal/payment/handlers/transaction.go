package handlers

import (
	"log"
	"net/http"

	"github.com/KrepkiyOrex/acquiring/internal/payment/service"
	"github.com/gorilla/mux"
)

func StartServer() {
	router := SetupRoutes()

	log.Println("===== Acquiring server started ... =====")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Listen and Server:", err)
	}
}

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", service.PaymentPage)
	router.HandleFunc("/pay", service.TransactionHandler)
	router.HandleFunc("/hey", service.Wel)

	return router
}
