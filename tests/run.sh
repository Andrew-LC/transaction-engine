#!/bin/bash

BASE_URL="http://localhost:8080"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
RESET='\033[0m'

PASS=0
FAIL=0

print_header() {
    echo ""
    echo -e "${CYAN}${BOLD}══════════════════════════════════════════${RESET}"
    echo -e "${CYAN}${BOLD}  $1${RESET}"
    echo -e "${CYAN}${BOLD}══════════════════════════════════════════${RESET}"
}

print_test() {
    echo ""
    echo -e "${YELLOW}▶ $1${RESET}"
}

run_curl() {
    local method=$1
    local path=$2
    local data=$3

    if [ -n "$data" ]; then
        curl -s -X "$method" "$BASE_URL$path" \
            -H "Content-Type: application/json" \
            -d "$data"
    else
        curl -s -X "$method" "$BASE_URL$path"
    fi
}

check() {
    local label=$1
    local response=$2
    local expected_field=$3
    local expected_value=$4

    actual=$(echo "$response" | grep -o "\"$expected_field\":\"[^\"]*\"" | cut -d'"' -f4)

    if [ "$actual" = "$expected_value" ]; then
        echo -e "  ${GREEN}✔ PASS${RESET} — $label"
        echo -e "  Response: $response"
        PASS=$((PASS + 1))
    else
        echo -e "  ${RED}✘ FAIL${RESET} — $label"
        echo -e "  Expected $expected_field=${BOLD}$expected_value${RESET}, got ${RED}$actual${RESET}"
        echo -e "  Response: $response"
        FAIL=$((FAIL + 1))
    fi
}

# ─────────────────────────────────────────────
print_header "1. CREATE CARD"
# ─────────────────────────────────────────────

print_test "Create a new card (Jane Smith)"
RESP=$(run_curl POST /api/card/newcard '{
  "card_number": 5000000000000001,
  "card_holder": "Jane Smith",
  "pin": "5678",
  "amount": 500
}')
check "New card created" "$RESP" "resp_code" "00"

print_test "Duplicate card (should fail)"
RESP=$(run_curl POST /api/card/newcard '{
  "card_number": 4123456789012345,
  "card_holder": "Duplicate",
  "pin": "0000",
  "amount": 0
}')
check "Duplicate card rejected" "$RESP" "status" "FAILED"

# ─────────────────────────────────────────────
print_header "2. GET BALANCE"
# ─────────────────────────────────────────────

print_test "Balance for seeded card (John Doe)"
RESP=$(run_curl GET /api/card/balance/4123456789012345)
check "Balance returns respCode 00" "$RESP" "resp_code" "00"

print_test "Balance for non-existent card"
RESP=$(run_curl GET /api/card/balance/9999999999999999)
check "Invalid card returns respCode 05" "$RESP" "resp_code" "05"

# ─────────────────────────────────────────────
print_header "3. WITHDRAW"
# ─────────────────────────────────────────────

print_test "Successful withdraw (200 from 1000)"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "withdraw",
  "amount": 200
}')
check "Withdraw succeeds with respCode 00" "$RESP" "resp_code" "00"

print_test "Insufficient balance withdraw"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "withdraw",
  "amount": 99999
}')
check "Insufficient balance returns respCode 99" "$RESP" "resp_code" "99"

print_test "Withdraw with wrong PIN"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "0000",
  "type": "withdraw",
  "amount": 100
}')
check "Wrong PIN returns respCode 06" "$RESP" "resp_code" "06"

print_test "Withdraw from non-existent card"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 9999999999999999,
  "pin": "1234",
  "type": "withdraw",
  "amount": 100
}')
check "Invalid card returns respCode 05" "$RESP" "resp_code" "05"

# ─────────────────────────────────────────────
print_header "4. TOPUP"
# ─────────────────────────────────────────────

print_test "Successful topup (500)"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "topup",
  "amount": 500
}')
check "Topup succeeds with respCode 00" "$RESP" "resp_code" "00"

print_test "Topup with wrong PIN"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "9999",
  "type": "topup",
  "amount": 100
}')
check "Wrong PIN on topup returns respCode 06" "$RESP" "resp_code" "06"

# ─────────────────────────────────────────────
print_header "5. INVALID TRANSACTION TYPE"
# ─────────────────────────────────────────────

print_test "Unknown transaction type"
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "transfer",
  "amount": 100
}')
check "Unknown type returns FAILED" "$RESP" "status" "FAILED"

# ─────────────────────────────────────────────
print_header "6. TRANSACTION HISTORY"
# ─────────────────────────────────────────────

print_test "Transaction history for John Doe"
RESP=$(run_curl GET /api/card/transactions/4123456789012345)
check "History returns respCode 00" "$RESP" "resp_code" "00"

print_test "Transaction history for non-existent card"
RESP=$(run_curl GET /api/card/transactions/9999999999999999)
check "Invalid card history returns respCode 05" "$RESP" "resp_code" "05"

# ─────────────────────────────────────────────
print_header "7. MALFORMED REQUEST"
# ─────────────────────────────────────────────

print_test "Malformed JSON body"
RESP=$(curl -s -X POST "$BASE_URL/api/transaction" \
    -H "Content-Type: application/json" \
    -d '{not valid json')
check "Malformed body returns FAILED" "$RESP" "status" "FAILED"

# ─────────────────────────────────────────────
print_header "8. END-TO-END FLOW"
# ─────────────────────────────────────────────

print_test "Create new card, topup, then withdraw"

run_curl POST /api/card/newcard '{
  "card_number": 5000000000000002,
  "card_holder": "Alice",
  "pin": "4321",
  "amount": 300
}' > /dev/null

run_curl POST /api/transaction '{
  "card_number": 5000000000000002,
  "pin": "4321",
  "type": "topup",
  "amount": 200
}' > /dev/null

RESP=$(run_curl POST /api/transaction '{
  "card_number": 5000000000000002,
  "pin": "4321",
  "type": "withdraw",
  "amount": 100
}')
check "E2E withdraw after topup succeeds" "$RESP" "resp_code" "00"

# ─────────────────────────────────────────────
echo ""
echo -e "${CYAN}${BOLD}══════════════════════════════════════════${RESET}"
echo -e "${BOLD}  RESULTS: ${GREEN}$PASS passed${RESET} / ${RED}$FAIL failed${RESET}"
echo -e "${CYAN}${BOLD}══════════════════════════════════════════${RESET}"
echo ""

if [ $FAIL -ne 0 ]; then
    exit 1
fi
