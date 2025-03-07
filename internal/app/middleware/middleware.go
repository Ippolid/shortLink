package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Секретный ключ для подписи JWT
var jwtSecret = []byte("super_secret_key")

// Структура JWT-токена
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Middleware аутентификации
//
//	func AuthMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			tokenStr, err := c.Cookie("jwt_token")
//			if err != nil {
//				// Если куки нет, создаём нового пользователя и JWT
//				newUserID := uuid.NewString()
//				tokenStr, err = createJWT(newUserID)
//				if err != nil {
//					c.AbortWithStatus(http.StatusInternalServerError)
//					return
//				}
//
//				// Устанавливаем куку с токеном
//				//c.SetCookie(c.Writer, &http.Cookie{
//				//	Name:     "jwt_token",
//				//	Value:    tokenStr,
//				//	HttpOnly: true,
//				//	Secure:   false, // Для HTTPS ставить true
//				//	Path:     "/",
//				//	Expires:  time.Now().Add(24 * time.Hour),
//				//})
//				c.SetCookie("jwt_token", tokenStr, 24*3600, "/", "", false, true)
//
//				c.Set("user_id", newUserID)
//			} else {
//				// Проверяем токен
//				claims, err := verifyJWT(tokenStr)
//				if err != nil {
//					c.AbortWithStatus(http.StatusUnauthorized)
//					return
//				}
//				c.Set("user_id", claims.UserID)
//			}
//
//			c.Next()
//		}
//	}
//
// // Функция создания JWT-токена
//
//	func createJWT(userID string) (string, error) {
//		claims := Claims{
//			UserID: userID,
//			RegisteredClaims: jwt.RegisteredClaims{
//				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 день
//				IssuedAt:  jwt.NewNumericDate(time.Now()),
//			},
//		}
//		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//		return token.SignedString(jwtSecret)
//	}
//
// // Функция проверки JWT-токена
//
//	func verifyJWT(tokenStr string) (*Claims, error) {
//		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//			return jwtSecret, nil
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//			return claims, nil
//		}
//		return nil, jwt.ErrSignatureInvalid
//	}
var secretKey = []byte("super_secret_key")

// Middleware аутентификации через куки
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user_session")

		if err != nil || !isValidCookie(cookie) {
			// Если куки нет или она недействительна — создаём новую
			newUserID := uuid.NewString()
			signedUserID := newUserID + ":" + signMessage(newUserID)

			// Устанавливаем куку с помощью Gin
			c.SetCookie("user_session", signedUserID, 24*3600, "/", "", false, true)

			// Передаём user_id в контекст запроса
			c.Set("user_id", newUserID)
		} else {
			// Извлекаем user_id из куки и передаём в контекст
			parts := splitCookie(cookie)
			c.Set("user_id", parts[0])
		}

		c.Next()
	}
}

// Функция проверки валидности куки
func isValidCookie(cookie string) bool {
	parts := splitCookie(cookie)
	if parts == nil || signMessage(parts[0]) != parts[1] {
		return false
	}
	return true
}

// Функция подписи строки (HMAC)
func signMessage(message string) string {
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Функция разделения куки на ID и подпись
func splitCookie(cookie string) []string {
	parts := make([]string, 2)
	for i, p := range []rune(cookie) {
		if p == ':' {
			parts[0] = cookie[:i]
			parts[1] = cookie[i+1:]
			return parts
		}
	}
	return nil
}
