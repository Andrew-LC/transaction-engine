# Transaction Engine

A backend service that processes card transactions and maintains card balances, simulating a simplified payment switch authorization engine. Built with Go using only standard libraries.

---

## Project Structure

```
transaction-engine/
├── cmd/              # Entry point
├── domain/           # Request/response types, constants
├── handler/          # HTTP handlers
├── middleware/       # Logger middleware
├── models/           # Card and Transaction structs
├── router/           # Route registration
├── service/          # Business logic
├── store/            # In-memory storage
└── tests/            # API integration tests
```

---

## Setup Instructions

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Docker](https://www.docker.com/get-started) (optional)

### Clone the Repository

```bash
git clone https://github.com/Andrew-LC/transaction-engine.git
cd transaction-engine
```

### Install Dependencies

```bash
go mod tidy
```

---

## Run Steps

### Option 1 — Go

```bash
go run ./cmd/...
```

### Option 2 — Docker

```bash
# Build the image
docker build -t transaction-engine .

# Run the container
docker run -p 8080:8080 transaction-engine
```

The server starts on `http://localhost:8080`.

A seeded card is available immediately on startup:

| Field       | Value                |
|-------------|----------------------|
| Card Number | `4123456789012345`   |
| Card Holder | John Doe             |
| PIN         | `1234`               |
| Balance     | `1000`               |
| Status      | `ACTIVE`             |

### Run Tests

```bash
go test ./tests/ -v
```

### Run curl Test Script

```bash
chmod +x run.sh
./run.sh
```

---

## API Reference

### Base URL

```
http://localhost:8080
```

### Endpoints

| Method | Endpoint                              | Description              |
|--------|---------------------------------------|--------------------------|
| POST   | `/api/card/newcard`                   | Create a new card        |
| POST   | `/api/transaction`                    | Process a transaction    |
| GET    | `/api/card/balance/{cardNumber}`      | Get card balance         |
| GET    | `/api/card/transactions/{cardNumber}` | Get transaction history  |

### Response Codes

| Code | Meaning              |
|------|----------------------|
| `00` | Success              |
| `05` | Invalid card         |
| `06` | Invalid PIN          |
| `99` | Insufficient balance |

---

## API Examples

### Create a New Card

```bash
curl -X POST http://localhost:8080/api/card/newcard \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": 5000000000000001,
    "card_holder": "Mikasa Ackerman",
    "pin": "5678",
    "amount": 500
  }'
```

**Response:**
```json
{
  "status": "SUCCESS",
  "resp_code": "00",
  "message": "created new card"
}
```

---

### Get Balance

```bash
curl http://localhost:8080/api/card/balance/4123456789012345
```

**Response:**
```json
{
  "status": "SUCCESS",
  "resp_code": "00",
  "balance": 1000
}
```

---

### Withdraw

```bash
curl -X POST http://localhost:8080/api/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": 4123456789012345,
    "pin": "1234",
    "type": "withdraw",
    "amount": 200
  }'
```

**Response:**
```json
{
  "status": "APPROVED",
  "resp_code": "00",
  "message": "transaction successful",
  "balance": 800
}
```

---

### Topup

```bash
curl -X POST http://localhost:8080/api/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": 4123456789012345,
    "pin": "1234",
    "type": "topup",
    "amount": 500
  }'
```

**Response:**
```json
{
  "status": "APPROVED",
  "resp_code": "00",
  "message": "transaction successful",
  "balance": 1500
}
```

---

### Get Transaction History

```bash
curl http://localhost:8080/api/card/transactions/4123456789012345
```

**Response:**
```json
{
  "status": "SUCCESS",
  "resp_code": "00",
  "transactions": [
    {
      "transactionId": "a3f2c1d4-...",
      "cardNumber": 4123456789012345,
      "type": "withdraw",
      "amount": 200,
      "status": "SUCCESS",
      "timestamp": "2026-03-25T13:36:13+05:30"
    }
  ]
}
```

---

### Error Responses

**Invalid Card (`05`):**
```bash
curl http://localhost:8080/api/card/balance/9999999999999999
```
```json
{
  "status": "FAILED",
  "resp_code": "05",
  "message": "Invalid card"
}
```

**Invalid PIN (`06`):**
```bash
curl -X POST http://localhost:8080/api/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": 4123456789012345,
    "pin": "0000",
    "type": "withdraw",
    "amount": 100
  }'
```
```json
{
  "status": "FAILED",
  "resp_code": "06",
  "message": "Invalid PIN"
}
```

**Insufficient Balance (`99`):**
```bash
curl -X POST http://localhost:8080/api/transaction \
  -H "Content-Type: application/json" \
  -d '{
    "card_number": 4123456789012345,
    "pin": "1234",
    "type": "withdraw",
    "amount": 99999
  }'
```
```json
{
  "status": "FAILED",
  "resp_code": "99",
  "message": "Insufficient balance"
}
```

---

## Postman

Import `postman_collection.json` from the repo root into Postman to run all requests with pre-built test assertions.

1. Open Postman → **Import**
2. Select `postman_collection.json`
3. Click **Run collection** to execute all 18 tests

---

## Security

- PINs are hashed using **SHA-256** before storage — plaintext PINs are never stored or logged
- All comparisons are done on hashed values only
- Concurrent access to the in-memory store is protected with `sync.RWMutex`
