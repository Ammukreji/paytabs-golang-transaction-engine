package models

import "time"

// Card represents a single user's card attributes
type Card struct {
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	PinHash    string  `json:"-"`
	Balance    float64 `json:"balance"`
	Status     string  `json:"status"` // ACTIVE or BLOCKED
}

// TransactionReq specifies the structure of a transaction request
type TransactionReq struct {
	CardNumber string  `json:"cardNumber"`
	Pin        string  `json:"pin"`
	Type       string  `json:"type"` // withdraw or topup
	Amount     float64 `json:"amount"`
}

// TransactionResp defines the structure for API responses regarding transactions
type TransactionResp struct {
	Status   string  `json:"status"`
	RespCode string  `json:"respCode"`
	Message  string  `json:"message,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
}

// TransactionLog represents a single log instance for a transaction
type TransactionLog struct {
	TransactionID string    `json:"transactionId"`
	CardNumber    string    `json:"cardNumber"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"` // SUCCESS or FAILED
	Timestamp     time.Time `json:"timestamp"`
}

// BalanceResp is structure representing the balance endpoint response
type BalanceResp struct {
	CardNumber string  `json:"cardNumber"`
	Balance    float64 `json:"balance"`
}
