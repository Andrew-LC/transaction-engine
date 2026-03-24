package router

import (
	"net/http"
	"transaction-engine/handler"
)

func Register(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/transaction", handler.ProcessTransaction)
	mux.HandleFunc("GET /api/card/balance/{cardNumber}", handler.GetBalance)
	mux.HandleFunc("GET /api/card/transactions/{cardNumber}", handler.GetTransactions)

	return mux
}
