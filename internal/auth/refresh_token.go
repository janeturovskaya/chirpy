package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	hexStr := hex.EncodeToString(b)
	return hexStr
}
