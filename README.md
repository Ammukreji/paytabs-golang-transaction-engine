# Transaction Processing Engine

A simplified backend service simulating a payment switch authorization engine. It processes card transactions and maintains card balances using a concurrent in-memory storage, built purely with Go's standard libraries.

## Setup Instructions

1. Ensure Go is installed on the machine.
2. Clone the repository and navigate to the project directory:
   ```bash
   cd paytabs-assessment
   ```
3. Initialize or verify standard modules:
   ```bash
   go mod tidy
   ```

## Run Steps

To run the application locally, execute the following command:
```bash
go run main.go
```

The terminal will log the following lines indicating a successful startup and database seed:
```text
Database seeded with sample card:
Card Number: 4123456789012345
Name: John Doe
PIN: 1234
Balance: 1000
Status: ACTIVE
-------------------------------------------------
Starting server on port 8080...
```

The service is now bound and ready to accept HTTP traffic at `http://localhost:8080`.

## API Examples (cURL)

**1. Process a Withdraw Transaction**
```bash
curl -X POST http://localhost:8080/api/transaction \
     -H "Content-Type: application/json" \
     -d '{
           "cardNumber": "4123456789012345",
           "pin": "1234",
           "type": "withdraw",
           "amount": 200
         }'
```

**2. Process a Topup Transaction**
```bash
curl -X POST http://localhost:8080/api/transaction \
     -H "Content-Type: application/json" \
     -d '{
           "cardNumber": "4123456789012345",
           "pin": "1234",
           "type": "topup",
           "amount": 500
         }'
```

**3. Check Current Account Balance**
```bash
curl http://localhost:8080/api/card/balance/4123456789012345
```

**4. View Historical Transactions**
```bash
curl http://localhost:8080/api/card/transactions/4123456789012345
```

## API Testing (Postman Collection)

To test the endpoints using Postman, copy the JSON configuration below and save it locally as a `postman_collection.json` file. Import this file into Postman to load all paths and payloads ready for execution.

```json
{
	"info": {
		"name": "Payment Switch API",
		"description": "API collection for the Transaction Processing Engine.",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [
		{
			"name": "Process Withdraw Transaction",
			"request": {
				"method": "POST",
				"header": [
					{
						"key": "Content-Type",
						"value": "application/json"
					}
				],
				"body": {
					"mode": "raw",
					"raw": "{\n  \"cardNumber\": \"4123456789012345\",\n  \"pin\": \"1234\",\n  \"type\": \"withdraw\",\n  \"amount\": 200\n}"
				},
				"url": {
					"raw": "http://localhost:8080/api/transaction",
					"protocol": "http",
					"host": ["localhost"],
					"port": "8080",
					"path": ["api", "transaction"]
				}
			}
		},
		{
			"name": "Process Topup Transaction",
			"request": {
				"method": "POST",
				"header": [
					{
						"key": "Content-Type",
						"value": "application/json"
					}
				],
				"body": {
					"mode": "raw",
					"raw": "{\n  \"cardNumber\": \"4123456789012345\",\n  \"pin\": \"1234\",\n  \"type\": \"topup\",\n  \"amount\": 500\n}"
				},
				"url": {
					"raw": "http://localhost:8080/api/transaction",
					"protocol": "http",
					"host": ["localhost"],
					"port": "8080",
					"path": ["api", "transaction"]
				}
			}
		},
		{
			"name": "Get Account Balance",
			"request": {
				"method": "GET",
				"url": {
					"raw": "http://localhost:8080/api/card/balance/4123456789012345",
					"protocol": "http",
					"host": ["localhost"],
					"port": "8080",
					"path": ["api", "card", "balance", "4123456789012345"]
				}
			}
		},
		{
			"name": "Get Transaction History",
			"request": {
				"method": "GET",
				"url": {
					"raw": "http://localhost:8080/api/card/transactions/4123456789012345",
					"protocol": "http",
					"host": ["localhost"],
					"port": "8080",
					"path": ["api", "card", "transactions", "4123456789012345"]
				}
			}
		}
	]
}
```
