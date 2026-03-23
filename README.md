# Transaction Processing Engine – Golang Backend Assessment

## Overview

This project is a backend service that processes card transactions and maintains card balances.
It simulates a simplified payment authorization engine using in-memory storage and SHA256 PIN authentication.

The system supports:

* Withdraw transactions
* Top-up transactions
* Balance inquiry
* Transaction history logging

---

## Technology Used

* Go (Golang)
* net/http (REST API)
* In-memory storage (map)
* SHA256 hashing
* UUID for transaction IDs

---

## Setup

Install dependencies:

```bash
go mod tidy
```

---

## Run

```bash
go run main.go
```

Server starts on:

```
http://localhost:8080
```

---

## Sample Cards

| Card Number      | PIN  | Balance | Status  |
| ---------------- | ---- | ------- | ------- |
| 4123456789012345 | 1234 | 1000    | ACTIVE  |
| 5123456789012345 | 5678 | 500     | ACTIVE  |
| 6123456789012345 | 9012 | 200     | BLOCKED |

---

## API Endpoints

### 1. Process Transaction

**POST** `/api/transaction`

```bash
curl -X POST http://localhost:8080/api/transaction \
  -H "Content-Type: application/json" \
  -d '{"cardNumber":"4123456789012345","pin":"1234","type":"withdraw","amount":200}'
```

Success Response:

```json
{"status":"SUCCESS","respCode":"00","balance":800}
```

Invalid Card:

```json
{"status":"FAILED","respCode":"05","message":"Invalid card"}
```

Invalid PIN:

```json
{"status":"FAILED","respCode":"06","message":"Invalid PIN"}
```

Insufficient Balance:

```json
{"status":"FAILED","respCode":"99","message":"Insufficient balance"}
```

---

### 2. Get Balance

```bash
curl http://localhost:8080/api/card/balance/4123456789012345
```

Response:

```json
{"cardNumber":"4123456789012345","balance":800}
```

---

### 3. Get Transaction History

```bash
curl http://localhost:8080/api/card/transactions/4123456789012345
```

Response:

```json
{
  "cardNumber":"4123456789012345",
  "transactions":[
    {
      "transactionId":"uuid",
      "cardNumber":"4123456789012345",
      "type":"withdraw",
      "amount":200,
      "status":"SUCCESS",
      "timestamp":"2024-01-15T10:30:00Z"
    }
  ]
}
```

---

## API Testing

### Using Curl

Examples are provided above for:

* Transaction
* Balance
* Transaction history

### Using Postman

Import the following JSON into Postman to test APIs:

```json
{
  "info": {
    "name": "Transaction Engine API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Process Transaction",
      "request": {
        "method": "POST",
        "header": [{"key": "Content-Type", "value": "application/json"}],
        "body": {
          "mode": "raw",
          "raw": "{\"cardNumber\":\"4123456789012345\",\"pin\":\"1234\",\"type\":\"withdraw\",\"amount\":200}"
        },
        "url": {"raw": "http://localhost:8080/api/transaction"}
      }
    },
    {
      "name": "Get Balance",
      "request": {
        "method": "GET",
        "url": {"raw": "http://localhost:8080/api/card/balance/4123456789012345"}
      }
    },
    {
      "name": "Get Transactions",
      "request": {
        "method": "GET",
        "url": {"raw": "http://localhost:8080/api/card/transactions/4123456789012345"}
      }
    }
  ]
}
```

---

## Response Codes

| Code | Description                 |
| ---- | --------------------------- |
| 00   | Success                     |
| 01   | Invalid request             |
| 05   | Invalid card / Card blocked |
| 06   | Invalid PIN                 |
| 07   | Invalid transaction type    |
| 08   | Invalid amount              |
| 99   | Insufficient balance        |

---

## Project Structure

```
.
├── main.go
├── go.mod
├── go.sum
├── README.md
```
