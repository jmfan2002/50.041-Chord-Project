package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// Commands for consistent hashing
func ConsistentHash(key string, nonce string) string {
	h := sha256.New()
	h.Write([]byte(key + nonce))
	b := h.Sum(nil)

	return hex.EncodeToString(b)
}
