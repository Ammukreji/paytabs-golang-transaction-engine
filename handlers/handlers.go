package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"paytabs-assessment/models"
	"paytabs-assessment/storage"
	"paytabs-assessment/utils"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// HandleTransaction processes /api/transaction endpoint
func HandleTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.TransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "90",
			Message:  "Invalid request format",
		})
		return
	}

	// Basic validation
	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "10",
			Message:  "Amount must be greater than 0",
		})
		return
	}
	if req.Type != "withdraw" && req.Type != "topup" {
		writeJSON(w, http.StatusBadRequest, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "11",
			Message:  "Invalid transaction type",
		})
		return
	}

	// Fetch Card
	card, exists := storage.DB.GetCard(req.CardNumber)
	if !exists {
		writeJSON(w, http.StatusOK, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Invalid card",
		})
		return
	}

	// Validate status
	if card.Status != "ACTIVE" {
		writeJSON(w, http.StatusOK, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "05",
			Message:  "Invalid card",
		})
		return
	}

	// Validate PIN securely
	if utils.HashPIN(req.Pin) != card.PinHash {
		writeJSON(w, http.StatusOK, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "06",
			Message:  "Invalid PIN",
		})
		return
	}

	txID := utils.GenerateTxID()

	// Process amount safely
	updatedCard, err := storage.DB.SubmitTransaction(req.CardNumber, req.Type, req.Amount)
	if err != nil {
		if err.Error() == "insufficient balance" {
			storage.DB.AppendTransactionHistory(models.TransactionLog{
				TransactionID: txID,
				CardNumber:    req.CardNumber,
				Type:          req.Type,
				Amount:        req.Amount,
				Status:        "FAILED",
				Timestamp:     time.Now(),
			})
			writeJSON(w, http.StatusOK, models.TransactionResp{
				Status:   "FAILED",
				RespCode: "99",
				Message:  "Insufficient balance",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, models.TransactionResp{
			Status:   "FAILED",
			RespCode: "91",
			Message:  err.Error(),
		})
		return
	}

	// Success logging
	storage.DB.AppendTransactionHistory(models.TransactionLog{
		TransactionID: txID,
		CardNumber:    req.CardNumber,
		Type:          req.Type,
		Amount:        req.Amount,
		Status:        "SUCCESS",
		Timestamp:     time.Now(),
	})

	writeJSON(w, http.StatusOK, models.TransactionResp{
		Status:   "SUCCESS",
		RespCode: "00",
		Balance:  updatedCard.Balance,
	})
}

// HandleCardEndpoints processes GET /api/card/balance/{cardNumber} and GET /api/card/transactions/{cardNumber}
func HandleCardEndpoints(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/card/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}

	action := parts[0]
	cardNumber := parts[1]

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	card, exists := storage.DB.GetCard(cardNumber)
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "Card not found",
		})
		return
	}

	if action == "balance" {
		writeJSON(w, http.StatusOK, models.BalanceResp{
			CardNumber: card.CardNumber,
			Balance:    card.Balance,
		})
		return
	} else if action == "transactions" {
		history := storage.DB.GetTransactionHistory(cardNumber)
		writeJSON(w, http.StatusOK, history)
		return
	}

	http.NotFound(w, r)
}
