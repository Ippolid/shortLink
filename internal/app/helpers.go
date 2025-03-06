package app

import (
	"crypto/sha1"
	"encoding/hex"
)

func GenerateShortID(url string, user string) string {
	hash := sha1.New()
	hash.Write([]byte(url + user))
	return hex.EncodeToString(hash.Sum(nil))[:8]
}
