package storage

import (
	"errors"
	"sync"

	"paytabs-assessment/models"
)

// MemoryDB simulates our mock database using Maps
type MemoryDB struct {
	cards        map[string]*models.Card
	transactions map[string][]models.TransactionLog
	mu           sync.RWMutex
}

var DB *MemoryDB

func InitDB() {
	DB = &MemoryDB{
		cards:        make(map[string]*models.Card),
		transactions: make(map[string][]models.TransactionLog),
	}
}

// AddCard inserts or overrides a single card
func (db *MemoryDB) AddCard(card models.Card) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.cards[card.CardNumber] = &card
}

// GetCard returns a copy of a particular card model if it exists
func (db *MemoryDB) GetCard(cardNumber string) (*models.Card, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	card, exists := db.cards[cardNumber]
	if !exists {
		return nil, false
	}
	cardCopy := *card
	return &cardCopy, true
}

// SubmitTransaction atomically processes a balance change on a card
func (db *MemoryDB) SubmitTransaction(cardNumber string, txType string, amount float64) (*models.Card, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	card, exists := db.cards[cardNumber]
	if !exists {
		return nil, errors.New("card not found")
	}

	if txType == "withdraw" {
		if card.Balance < amount {
			return nil, errors.New("insufficient balance")
		}
		card.Balance -= amount
	} else if txType == "topup" {
		card.Balance += amount
	} else {
		return nil, errors.New("invalid transaction type")
	}

	cardCopy := *card
	return &cardCopy, nil
}

// AppendTransactionHistory logs a successful or failed action
func (db *MemoryDB) AppendTransactionHistory(log models.TransactionLog) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.transactions[log.CardNumber] = append(db.transactions[log.CardNumber], log)
}

// GetTransactionHistory retrieves recent log details
func (db *MemoryDB) GetTransactionHistory(cardNumber string) []models.TransactionLog {
	db.mu.RLock()
	defer db.mu.RUnlock()
	txs, exists := db.transactions[cardNumber]
	if !exists {
		return []models.TransactionLog{} // return empty slice
	}
	
	txsCopy := make([]models.TransactionLog, len(txs))
	copy(txsCopy, txs)
	return txsCopy
}
