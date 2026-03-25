package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"transaction-engine/domain"
	"transaction-engine/service"
	"transaction-engine/store"
)

func parseCardNumber(w http.ResponseWriter, r *http.Request) (int64, bool) {
	str := r.PathValue("cardNumber")
	if str == "" {
		writeJSON(w, http.StatusBadRequest, domain.NewResponse(
			"FAILED", "03", "missing card number", 0,
		))
		return 0, false
	}
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, domain.NewResponse(
			"FAILED", "03", "invalid card number", 0,
		))
		return 0, false
	}
	return n, true
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrCardNotFound):
		writeJSON(w, http.StatusNotFound, domain.NewResponse(
			"FAILED", "05", "invalid card", 0,
		))
	case errors.Is(err, store.ErrCardExists):
		writeJSON(w, http.StatusNotFound, domain.NewResponse(
			"FAILED", "05", "card already exists", 0,
		))
	case errors.Is(err, service.ErrInvalidPin):
		writeJSON(w, http.StatusUnauthorized, domain.NewResponse(
			"FAILED", "06", "invalid pin", 0,
		))
	case errors.Is(err, service.ErrInsufficientFunds):
		writeJSON(w, http.StatusBadRequest, domain.NewResponse(
			"FAILED", "99", "insufficient funds", 0,
		))
	case errors.Is(err, service.ErrInvalidTransactionType):
		writeJSON(w, http.StatusBadRequest, domain.NewResponse(
			"FAILED", "12", "invalid transaction type", 0,
		))
	default:
		writeJSON(w, http.StatusInternalServerError, domain.NewResponse(
			"FAILED", "99", "internal error", 0,
		))
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
