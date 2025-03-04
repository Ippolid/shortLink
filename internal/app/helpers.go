package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/models"
)

func GenerateShortID(url []byte) string {
	hash := sha1.New()
	hash.Write(url)
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

func SignUserID(userID string) string {
	h := hmac.New(sha256.New, []byte(config.SecretKey))
	h.Write([]byte(userID))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateCookie(value string) (string, bool) {
	var data models.UserCookie

	err := json.Unmarshal([]byte(value), &data)

	fmt.Println(data)
	if err != nil {
		return "", false
	}

	expectedSign := SignUserID(data.UserID)
	if data.Sign != expectedSign {
		return "", false
	}

	return data.UserID, true
}
