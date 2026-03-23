package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateTxID uses crypto/rand to give a randomized hex string id.
func GenerateTxID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "tx-error" // fallback
	}
	return hex.EncodeToString(b)
}
