package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CardStatus string

const (
	StatusActive  CardStatus = "ACTIVE"
	StatusBlocked CardStatus = "BLOCKED"
)

type TransactionType string

const (
	TypeWithdraw TransactionType = "withdraw"
	TypeTopup    TransactionType = "topup"
)

type TransactionStatus string

const (
	TxSuccess TransactionStatus = "SUCCESS"
	TxFailed  TransactionStatus = "FAILED"
)

type Card struct {
	CardNumber string     `json:"cardNumber"`
	CardHolder string     `json:"cardHolder"`
	PinHash    string     `json:"-"`
	Balance    float64    `json:"balance"`
	Status     CardStatus `json:"status"`
}

type Transaction struct {
	TransactionID string            `json:"transactionId"`
	CardNumber    string            `json:"cardNumber"`
	Type          TransactionType  `json:"type"`
	Amount        float64           `json:"amount"`
	Status        TransactionStatus `json:"status"`
	Timestamp     time.Time         `json:"timestamp"`
}

type TransactionRequest struct {
	CardNumber string  `json:"cardNumber"`
	Pin        string  `json:"pin"`
	Type       string  `json:"type"`
	Amount     float64 `json:"amount"`
}

type TransactionResponse struct {
	Status   string  `json:"status"`
	RespCode string  `json:"respCode,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
	Message  string  `json:"message,omitempty"`
}

type CardService struct {
	cards map[string]*Card
	mu    sync.RWMutex
}

func NewCardService() *CardService {
	service := &CardService{
		cards: make(map[string]*Card),
	}
	service.initSampleCards()
	return service
}

func (s *CardService) initSampleCards() {
	s.cards["4123456789012345"] = &Card{
		CardNumber: "4123456789012345",
		CardHolder: "John Doe",
		PinHash:    HashPin("1234"),
		Balance:    1000,
		Status:     StatusActive,
	}
	s.cards["5123456789012345"] = &Card{
		CardNumber: "5123456789012345",
		CardHolder: "Jane Smith",
		PinHash:    HashPin("5678"),
		Balance:    500,
		Status:     StatusActive,
	}
	s.cards["6123456789012345"] = &Card{
		CardNumber: "6123456789012345",
		CardHolder: "Bob Wilson",
		PinHash:    HashPin("9012"),
		Balance:    200,
		Status:     StatusBlocked,
	}
}

func HashPin(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}

func (s *CardService) GetCard(cardNumber string) (*Card, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, exists := s.cards[cardNumber]
	return card, exists
}

func (s *CardService) GetBalance(cardNumber string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, exists := s.cards[cardNumber]
	if !exists {
		return 0, false
	}
	return card.Balance, true
}

func (s *CardService) UpdateBalance(cardNumber string, newBalance float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	card, exists := s.cards[cardNumber]
	if !exists {
		return false
	}
	card.Balance = newBalance
	return true
}

type TransactionService struct {
	transactions []Transaction
	cardService  *CardService
	mu           sync.Mutex
}

func NewTransactionService(cardService *CardService) *TransactionService {
	return &TransactionService{
		transactions: make([]Transaction, 0),
		cardService:  cardService,
	}
}

func (s *TransactionService) ProcessTransaction(req TransactionRequest) TransactionResponse {
	card, exists := s.cardService.GetCard(req.CardNumber)
	if !exists {
		return TransactionResponse{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Invalid card",
		}
	}

	if card.Status != StatusActive {
		return TransactionResponse{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Card is blocked",
		}
	}

	inputPinHash := HashPin(req.Pin)
	if inputPinHash != card.PinHash {
		return TransactionResponse{
			Status:   "FAILED",
			RespCode: "06",
			Message:  "Invalid PIN",
		}
	}

	if req.Type != string(TypeWithdraw) && req.Type != string(TypeTopup) {
		return TransactionResponse{
			Status:   "FAILED",
			RespCode: "07",
			Message:  "Invalid transaction type",
		}
	}

	if req.Amount <= 0 {
		return TransactionResponse{
			Status:   "FAILED",
			RespCode: "08",
			Message:  "Invalid amount",
		}
	}

	txType := TransactionType(req.Type)
	var newBalance float64
	var status TransactionStatus

	if txType == TypeWithdraw {
		if card.Balance < req.Amount {
			s.logTransaction(req, TxFailed)
			return TransactionResponse{
				Status:   "FAILED",
				RespCode: "99",
				Message:  "Insufficient balance",
			}
		}
		newBalance = card.Balance - req.Amount
		status = TxSuccess
	} else {
		newBalance = card.Balance + req.Amount
		status = TxSuccess
	}

	s.cardService.UpdateBalance(req.CardNumber, newBalance)
	s.logTransaction(req, status)

	return TransactionResponse{
		Status:   "SUCCESS",
		RespCode: "00",
		Balance:  newBalance,
	}
}

func (s *TransactionService) logTransaction(req TransactionRequest, status TransactionStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx := Transaction{
		TransactionID: uuid.New().String(),
		CardNumber:    req.CardNumber,
		Type:          TransactionType(req.Type),
		Amount:        req.Amount,
		Status:        status,
		Timestamp:     time.Now(),
	}
	s.transactions = append(s.transactions, tx)
}

func (s *TransactionService) GetTransactions(cardNumber string) []Transaction {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Transaction, 0)
	for _, tx := range s.transactions {
		if tx.CardNumber == cardNumber {
			result = append(result, tx)
		}
	}
	return result
}

type Handler struct {
	cardService        *CardService
	transactionService *TransactionService
}

func NewHandler() *Handler {
	cardService := NewCardService()
	transactionService := NewTransactionService(cardService)
	return &Handler{
		cardService:        cardService,
		transactionService: transactionService,
	}
}

func (h *Handler) handleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, TransactionResponse{
			Status:   "FAILED",
			RespCode: "01",
			Message:  "Invalid request",
		}, http.StatusBadRequest)
		return
	}

	response := h.transactionService.ProcessTransaction(req)
	sendJSON(w, response, http.StatusOK)
}

func (h *Handler) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := extractCardNumber(r.URL.Path)
	balance, exists := h.cardService.GetBalance(cardNumber)
	if !exists {
		sendJSON(w, map[string]string{
			"status":  "FAILED",
			"message": "Card not found",
		}, http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"cardNumber": cardNumber,
		"balance":    balance,
	}, http.StatusOK)
}

func (h *Handler) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cardNumber := extractCardNumber(r.URL.Path)
	transactions := h.transactionService.GetTransactions(cardNumber)

	sendJSON(w, map[string]interface{}{
		"cardNumber":   cardNumber,
		"transactions": transactions,
	}, http.StatusOK)
}

func extractCardNumber(path string) string {
	if strings.HasPrefix(path, "/api/card/balance/") {
		return strings.TrimPrefix(path, "/api/card/balance/")
	}
	if strings.HasPrefix(path, "/api/card/transactions/") {
		return strings.TrimPrefix(path, "/api/card/transactions/")
	}
	return ""
}

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	handler := NewHandler()

	http.HandleFunc("/api/transaction", handler.handleTransaction)
	http.HandleFunc("/api/card/balance/", handler.handleGetBalance)
	http.HandleFunc("/api/card/transactions/", handler.handleGetTransactions)

	fmt.Println("Transaction Processing Engine starting on :8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  POST   /api/transaction")
	fmt.Println("  GET    /api/card/balance/{cardNumber}")
	fmt.Println("  GET    /api/card/transactions/{cardNumber}")
	fmt.Println()
	fmt.Println("Sample cards:")
	fmt.Println("  Card: 4123456789012345, PIN: 1234, Balance: 1000, Status: ACTIVE")
	fmt.Println("  Card: 5123456789012345, PIN: 5678, Balance: 500, Status: ACTIVE")
	fmt.Println("  Card: 6123456789012345, PIN: 9012, Balance: 200, Status: BLOCKED")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
