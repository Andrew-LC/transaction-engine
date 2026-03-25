package utils

import (
	"encoding/hex"
	"crypto/sha256"
)

func Hash(pin string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pin))

	return hex.EncodeToString(hasher.Sum(nil)) 
}
