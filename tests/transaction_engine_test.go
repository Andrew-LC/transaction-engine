package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"transaction-engine/domain"
	"transaction-engine/handler"
	"transaction-engine/router"
	"transaction-engine/service"
	"transaction-engine/store"
)


func setupServer() http.Handler {
	st := store.NewStore() 
	svc := service.NewService(st)
	h := handler.NewHandler(svc)
	return router.Register(h)
}

func doRequest(t *testing.T, mux http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder) domain.Response {
	t.Helper()
	var resp domain.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	return resp
}

const (
	seededCard   = 4123456789012345
	seededPin    = "1234"
	seededBal    = 1000
)


func TestWithdraw_Success(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "withdraw",
		"amount":      200,
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00, got %s", resp.RespCode)
	}
	if resp.Balance != 800 {
		t.Errorf("expected balance 800, got %d", resp.Balance)
	}
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "withdraw",
		"amount":      9999,
	})

	resp := decodeResponse(t, rr)
	if resp.RespCode != "99" {
		t.Errorf("expected respCode 99, got %s", resp.RespCode)
	}
	if resp.Status != "FAILED" {
		t.Errorf("expected status FAILED, got %s", resp.Status)
	}
}

func TestWithdraw_ExactBalance(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "withdraw",
		"amount":      seededBal,
	})

	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00 for exact-balance withdraw, got %s", resp.RespCode)
	}
	if resp.Balance != 0 {
		t.Errorf("expected balance 0, got %d", resp.Balance)
	}
}


func TestTopup_Success(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "topup",
		"amount":      500,
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00, got %s", resp.RespCode)
	}
	if resp.Balance != 1500 {
		t.Errorf("expected balance 1500, got %d", resp.Balance)
	}
}

// POST /api/transaction — invalid card
func TestTransaction_InvalidCard(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": int64(9999999999999999),
		"pin":         seededPin,
		"type":        "withdraw",
		"amount":      100,
	})

	resp := decodeResponse(t, rr)
	if resp.RespCode != "05" {
		t.Errorf("expected respCode 05 for invalid card, got %s", resp.RespCode)
	}
	if resp.Status != "FAILED" {
		t.Errorf("expected status FAILED, got %s", resp.Status)
	}
}

// POST /api/transaction — invalid PIN
func TestTransaction_InvalidPin(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         "0000",
		"type":        "withdraw",
		"amount":      100,
	})

	resp := decodeResponse(t, rr)
	if resp.RespCode != "06" {
		t.Errorf("expected respCode 06 for invalid PIN, got %s", resp.RespCode)
	}
	if resp.Status != "FAILED" {
		t.Errorf("expected status FAILED, got %s", resp.Status)
	}
}

// POST /api/transaction — invalid transaction type
func TestTransaction_InvalidType(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "transfer",
		"amount":      100,
	})

	resp := decodeResponse(t, rr)
	if resp.Status != "FAILED" {
		t.Errorf("expected FAILED for unknown tx type, got %s", resp.Status)
	}
}

// POST /api/transaction — bad request body
func TestTransaction_MalformedBody(t *testing.T) {
	mux := setupServer()

	req := httptest.NewRequest(http.MethodPost, "/api/transaction", bytes.NewBufferString("{not valid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed body, got %d", rr.Code)
	}
}

// GET /api/card/balance/{cardNumber}
func TestGetBalance_Success(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodGet, "/api/card/balance/4123456789012345", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00, got %s", resp.RespCode)
	}
	if resp.Balance != seededBal {
		t.Errorf("expected balance %d, got %d", seededBal, resp.Balance)
	}
}

func TestGetBalance_InvalidCard(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodGet, "/api/card/balance/9999999999999999", nil)

	resp := decodeResponse(t, rr)
	if resp.RespCode != "05" {
		t.Errorf("expected respCode 05 for unknown card, got %s", resp.RespCode)
	}
}

func TestGetBalance_ReflectsAfterWithdraw(t *testing.T) {
	// Balance endpoint must return updated balance after a transaction
	mux := setupServer()

	doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard,
		"pin":         seededPin,
		"type":        "withdraw",
		"amount":      300,
	})

	rr := doRequest(t, mux, http.MethodGet, "/api/card/balance/4123456789012345", nil)
	resp := decodeResponse(t, rr)

	if resp.Balance != 700 {
		t.Errorf("expected balance 700 after withdraw, got %d", resp.Balance)
	}
}

// GET /api/card/transactions/{cardNumber}
func TestGetTransactions_EmptyInitially(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodGet, "/api/card/transactions/4123456789012345", nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00, got %s", resp.RespCode)
	}
}

func TestGetTransactions_AppearsAfterTransaction(t *testing.T) {
	mux := setupServer()

	doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard, "pin": seededPin, "type": "withdraw", "amount": 100,
	})
	doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": seededCard, "pin": seededPin, "type": "topup", "amount": 50,
	})

	rr := doRequest(t, mux, http.MethodGet, "/api/card/transactions/4123456789012345", nil)

	var raw map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	txs, ok := raw["transactions"].([]any)
	if !ok || len(txs) != 2 {
		t.Errorf("expected 2 transactions, got %v", raw["transactions"])
	}
}

func TestGetTransactions_InvalidCard(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodGet, "/api/card/transactions/9999999999999999", nil)

	resp := decodeResponse(t, rr)
	if resp.RespCode != "05" {
		t.Errorf("expected respCode 05 for unknown card, got %s", resp.RespCode)
	}
}

// POST /api/card/newcard
func TestCreateCard_Success(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/card/newcard", map[string]any{
		"card_number": int64(5000000000000001),
		"card_holder": "Jane Smith",
		"pin":         "5678",
		"amount":      500,
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected respCode 00, got %s", resp.RespCode)
	}
}

func TestCreateCard_DuplicateCard(t *testing.T) {
	mux := setupServer()

	rr := doRequest(t, mux, http.MethodPost, "/api/card/newcard", map[string]any{
		"card_number": seededCard,
		"card_holder": "Duplicate",
		"pin":         "0000",
		"amount":      0,
	})

	resp := decodeResponse(t, rr)
	if resp.Status != "FAILED" {
		t.Errorf("expected FAILED for duplicate card, got %s", resp.Status)
	}
}

func TestCreateCard_ThenTransact(t *testing.T) {
	mux := setupServer()

	doRequest(t, mux, http.MethodPost, "/api/card/newcard", map[string]any{
		"card_number": int64(5000000000000002),
		"card_holder": "Alice",
		"pin":         "4321",
		"amount":      300,
	})

	rr := doRequest(t, mux, http.MethodPost, "/api/transaction", map[string]any{
		"card_number": int64(5000000000000002),
		"pin":         "4321",
		"type":        "withdraw",
		"amount":      100,
	})

	resp := decodeResponse(t, rr)
	if resp.RespCode != "00" {
		t.Errorf("expected 00 after withdraw on new card, got %s", resp.RespCode)
	}
	if resp.Balance != 200 {
		t.Errorf("expected balance 200, got %d", resp.Balance)
	}
}
