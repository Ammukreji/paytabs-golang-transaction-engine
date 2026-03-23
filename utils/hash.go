package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPIN securely hashes a given pin string using SHA-256
func HashPIN(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return hex.EncodeToString(hash[:])
}
