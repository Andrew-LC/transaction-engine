package handler

import (
	"encoding/json"
	"net/http"

	"transaction-engine/domain"
	"transaction-engine/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var request domain.NewCardRequest
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, domain.NewResponse(
			"FAILED", "03", "invalid request body", 0,
		))
		return
	}
	
	_, err := h.service.CreateCard(request)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, domain.NewResponse(
		"SUCCESS", "00", "created new card", 0,
	))
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
    cardNumber, ok := parseCardNumber(w, r)
    if !ok {
        return
    }

    balance, err := h.service.GetBalance(cardNumber)
    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, domain.NewResponse(
        "SUCCESS", "00", "", balance,
    ))
}

func (h *Handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
    cardNumber, ok := parseCardNumber(w, r)
    if !ok {
        return
    }

    transactions, err := h.service.GetTransactions(cardNumber)
    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, domain.NewTransactionsResponse(transactions))
}

func (h *Handler) ProcessTransaction(w http.ResponseWriter, r *http.Request) {
    var req domain.TransactionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, domain.NewResponse(
            "FAILED", "03", "invalid request body", 0,
        ))
        return
    }

    newBalance, err := h.service.ProcessTransaction(req)
    if err != nil {
        writeError(w, err)
        return
    }

    writeJSON(w, http.StatusOK, domain.NewResponse(
        "APPROVED", "00", "transaction successful", newBalance,
    ))
}
