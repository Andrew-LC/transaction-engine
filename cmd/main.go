package main

import (
	"log"
	"net/http"

	"transaction-engine/handler"
	"transaction-engine/service"
	"transaction-engine/store"
)


func main() {
	mux := http.NewServeMux()

	s := store.NewStore()
	svc := service.NewService(s)
	handler := handler.NewHandler(svc)
	
	mux.HandleFunc("POST /api/transaction", handler.ProcessTransaction)
	mux.HandleFunc("GET /api/card/balance/{cardNumber}", handler.GetBalance)
	mux.HandleFunc("GET /api/card/transactions/{cardNumber}", handler.GetTransactions)
	
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
