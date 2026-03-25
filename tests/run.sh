#!/bin/bash

BASE_URL="http://localhost:8080"
PASS=0
FAIL=0
RESULTS=()

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
    local section=$1
    local label=$2
    local response=$3
    local expected_field=$4
    local expected_value=$5

    actual=$(echo "$response" | grep -o "\"$expected_field\":\"[^\"]*\"" | cut -d'"' -f4)

    if [ "$actual" = "$expected_value" ]; then
        RESULTS+=("$(printf '%-30s | %-38s | %-6s | PASS' "$section" "$label" "$expected_value")")
        PASS=$((PASS + 1))
    else
        RESULTS+=("$(printf '%-30s | %-38s | %-6s | FAIL  (got: %s)' "$section" "$label" "$expected_value" "$actual")")
        FAIL=$((FAIL + 1))
    fi
}

divider() {
    printf '%s\n' "------------------------------+----------------------------------------+--------+----------------------"
}

# 1. Create Card
RESP=$(run_curl POST /api/card/newcard '{
  "card_number": 5000000000000001,
  "card_holder": "Mikasa Ackerman",
  "pin": "5678",
  "amount": 500
}')
check "1. Create Card" "New card created" "$RESP" "resp_code" "00"

RESP=$(run_curl POST /api/card/newcard '{
  "card_number": 4123456789012345,
  "card_holder": "Duplicate",
  "pin": "0000",
  "amount": 0
}')
check "1. Create Card" "Duplicate card rejected" "$RESP" "status" "FAILED"

# 2. Get Balance
RESP=$(run_curl GET /api/card/balance/4123456789012345)
check "2. Get Balance" "Seeded card balance" "$RESP" "resp_code" "00"

RESP=$(run_curl GET /api/card/balance/9999999999999999)
check "2. Get Balance" "Non-existent card" "$RESP" "resp_code" "05"

# 3. Withdraw
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "withdraw",
  "amount": 200
}')
check "3. Withdraw" "Success (200 from 1000)" "$RESP" "resp_code" "00"

RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "withdraw",
  "amount": 99999
}')
check "3. Withdraw" "Insufficient balance" "$RESP" "resp_code" "99"

RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "0000",
  "type": "withdraw",
  "amount": 100
}')
check "3. Withdraw" "Wrong PIN" "$RESP" "resp_code" "06"

RESP=$(run_curl POST /api/transaction '{
  "card_number": 9999999999999999,
  "pin": "1234",
  "type": "withdraw",
  "amount": 100
}')
check "3. Withdraw" "Non-existent card" "$RESP" "resp_code" "05"

# 4. Topup
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "topup",
  "amount": 500
}')
check "4. Topup" "Success (500)" "$RESP" "resp_code" "00"

RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "9999",
  "type": "topup",
  "amount": 100
}')
check "4. Topup" "Wrong PIN" "$RESP" "resp_code" "06"

RESP=$(run_curl POST /api/transaction '{
  "card_number": 9999999999999999,
  "pin": "1234",
  "type": "topup",
  "amount": 100
}')
check "4. Topup" "Non-existent card" "$RESP" "resp_code" "05"

# 5. Invalid Transaction Type
RESP=$(run_curl POST /api/transaction '{
  "card_number": 4123456789012345,
  "pin": "1234",
  "type": "transfer",
  "amount": 100
}')
check "5. Invalid Type" "Unknown type (transfer)" "$RESP" "status" "FAILED"

RESP=$(curl -s -X POST "$BASE_URL/api/transaction" \
    -H "Content-Type: application/json" \
    -d '{not valid json')
check "5. Invalid Type" "Malformed JSON body" "$RESP" "status" "FAILED"

# 6. Transaction History
RESP=$(run_curl GET /api/card/transactions/4123456789012345)
check "6. Tx History" "Seeded card history" "$RESP" "resp_code" "00"

RESP=$(run_curl GET /api/card/transactions/9999999999999999)
check "6. Tx History" "Non-existent card" "$RESP" "resp_code" "05"

# 7. End-to-End Flow
run_curl POST /api/card/newcard '{
  "card_number": 5000000000000002,
  "card_holder": "Rem",
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
check "7. End-to-End" "Create -> topup -> withdraw" "$RESP" "resp_code" "00"


echo ""
echo "Transaction Engine — Test Results"
echo "Run at: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""
printf '%-30s | %-38s | %-6s | %s\n' "Section" "Test" "Expect" "Result"
divider
for row in "${RESULTS[@]}"; do
    echo "$row"
done
divider
echo ""
printf 'Total: %d tests    Passed: %d    Failed: %d\n' "$((PASS + FAIL))" "$PASS" "$FAIL"
echo ""

if [ $FAIL -ne 0 ]; then
    exit 1
fi
