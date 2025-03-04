package app

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Ippolid/shortLink/config"
	"github.com/gin-gonic/gin"
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

func GetUserId(c *gin.Context) (string, error) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return "", fmt.Errorf("user id not found")
	}

	return userIDVal.(string), nil
}
