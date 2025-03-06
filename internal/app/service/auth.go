package service

import (
	"fmt"
	"github.com/Ippolid/shortLink/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var SecretSalt = []byte("practicumSecretKey32")
var tokenSalt = []byte("tokenPracticum32")

func VerifyUser(token string) (string, error) {
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return SecretSalt, nil
	})
	if err != nil || !parsedToken.Valid {
		return "", fmt.Errorf("incorrect token: %v", err)
	}

	return claims.UserID, nil
}

func CreatTokenForUser(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	signedToken, err := token.SignedString(SecretSalt)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
