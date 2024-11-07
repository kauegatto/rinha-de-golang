package main

import (
	"log"
	"net/http"
	"rinha_backend/cmd/api/handlers"
	"rinha_backend/internal/domain/services"
	"rinha_backend/internal/infrastructure"
	"rinha_backend/internal/infrastructure/postgres"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	db, err := infrastructure.NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	clientRepo := postgres.NewClientRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)

	clientService := services.NewClientService(clientRepo, transactionRepo)
	clientHandler := handlers.NewClientHandler(clientService)

	http.HandleFunc("POST /clientes/{id}/transacoes", clientHandler.HandleTransaction)
	http.HandleFunc("GET /clientes/{id}/extrato", clientHandler.HandleExtract)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
