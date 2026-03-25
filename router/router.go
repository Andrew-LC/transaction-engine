package router

import (
	"net/http"
	"transaction-engine/handler"
	"transaction-engine/middleware"
)

func Register(handler *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	middlewares := []func(http.Handler) http.Handler {
		middleware.Logger,
	}

	mux.Handle("POST /api/transaction", middleware.Chain(http.HandlerFunc(handler.ProcessTransaction), middlewares...))
	mux.Handle("GET /api/card/balance/{cardNumber}", middleware.Chain(http.HandlerFunc(handler.GetBalance), middlewares...))
	mux.Handle("GET /api/card/transactions/{cardNumber}", middleware.Chain(http.HandlerFunc(handler.GetTransactions), middlewares...))

	return mux
}
