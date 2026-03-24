package main

import (
	"log"
	"net/http"

	"transaction-engine/handler"
	"transaction-engine/service"
	"transaction-engine/store"
	"transaction-engine/router"
)


func main() {
	s := store.NewStore()
	svc := service.NewService(s)
	handler := handler.NewHandler(svc)
	mux := router.Register(handler)
	
	
	log.Println("Server started in :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
