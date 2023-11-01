package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256String(val string) string {
	hash := sha256.Sum256([]byte(val))
	return hex.EncodeToString(hash[:])
}
