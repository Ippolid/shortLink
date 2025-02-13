package app

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateShortID(url []byte) string {
	hash := sha1.New()
	hash.Write(url)
	return hex.EncodeToString(hash.Sum(nil))[:8]
}
