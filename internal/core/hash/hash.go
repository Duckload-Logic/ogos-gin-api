package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSHA256Hash(input string, size int) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:size])
}
