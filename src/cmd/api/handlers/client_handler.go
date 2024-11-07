package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rinha_backend/internal/domain/models"
	"rinha_backend/internal/domain/services"
)

type ClientHandler struct {
	clientService *services.ClientService
}

func NewClientHandler(cs *services.ClientService) *ClientHandler {
	return &ClientHandler{clientService: cs}
}

func (h *ClientHandler) HandleTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	client, err := h.clientService.ProcessTransaction(r.Context(), id, payment)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"limite": client.AccountLimit,
		"saldo":  client.Balance,
	})
}

func (h *ClientHandler) HandleExtract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	client, transactions, err := h.clientService.GetExtract(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	response := struct {
		AccountBalance models.Client        `json:"saldo"`
		Transactions   []models.Transaction `json:"ultimas_transacoes"`
	}{
		AccountBalance: client,
		Transactions:   transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
