package middleware

import (
	"github.com/Ippolid/shortLink/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"
)

//	func AuthMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			authHeader := c.GetHeader("Authorization")
//			var accessToken string
//			if authHeader != "" {
//				accessToken = authHeader
//			} else {
//				cookie, err := c.Cookie("user_id")
//				if err == nil && cookie != "" {
//					accessToken = cookie
//				} else {
//					accessToken = ""
//				}
//			}
//
//			userID, err := service.VerifyUser(accessToken)
//			if err != nil {
//				userID = uuid.New().String()
//				token, err := service.CreatTokenForUser(userID)
//				if err != nil {
//					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
//					c.Abort()
//					return
//				}
//
//				c.SetCookie("user_id", token, 3600, "/", "", false, true)
//				c.Header("Authorization", token)
//			}
//
//			c.Set("user_id", userID)
//			c.Next()
//		}
//	}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var accessToken string
		if authHeader != "" {
			accessToken = authHeader
		} else {
			cookie, err := c.Cookie("user_id")
			if err == nil && cookie != "" {
				accessToken = cookie
			} else {
				accessToken = ""
			}
		}

		userID, err := service.VerifyUser(accessToken)
		if err != nil || accessToken == "" {
			userID = uuid.New().String()
			token, err := service.CreatTokenForUser(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
				c.Abort()
				return
			}

			c.SetCookie("user_id", token, 3600, "/", "", false, true)
			c.Header("Authorization", token)
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

//func AuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.GetHeader("Authorization")
//		var accessToken string
//		if authHeader != "" {
//			accessToken = authHeader
//		} else {
//			cookie, err := c.Cookie("user_id")
//			if err == nil && cookie != "" {
//				accessToken = cookie
//			}
//		}
//
//		userID, err := service.VerifyUser(accessToken)
//		if err != nil || accessToken == "" {
//			// Создаём новый идентификатор и токен
//			userID = uuid.New().String()
//			token, err := service.CreatTokenForUser(userID)
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
//				c.Abort()
//				return
//			}
//			c.SetCookie("user_id", token, 3600, "/", "", false, true)
//			c.Header("Authorization", token)
//		}
//
//		c.Set("user_id", userID)
//		c.Next()
//	}
//}
//func CheckAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		accessToken, err := c.Cookie("user_id")
//		if err != nil {
//			c.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		userID, err := service.VerifyUser(accessToken)
//		if err != nil {
//			c.AbortWithStatus(http.StatusUnauthorized)
//			return
//		}
//
//		c.Set("user_id", userID)
//		c.Next()
//	}
//}

// CheckAuthMiddleware теперь не возвращает 401, а создаёт новый токен при его отсутствии.
func CheckAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("user_id")
		if err != nil || accessToken == "" {
			// Создание нового токена
			userID := uuid.New().String()
			token, err := service.CreatTokenForUser(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
				c.Abort()
				return
			}
			c.SetCookie("user_id", token, 3600, "/", "", false, true)
			c.Header("Authorization", token)
			c.Set("user_id", userID)
		} else {
			userID, err := service.VerifyUser(accessToken)
			if err != nil {
				// Если верификация не прошла, создаём новый токен
				userID = uuid.New().String()
				token, err := service.CreatTokenForUser(userID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate auth token"})
					c.Abort()
					return
				}
				c.SetCookie("user_id", token, 3600, "/", "", false, true)
				c.Header("Authorization", token)
			}
			c.Set("user_id", userID)
		}
		c.Next()
	}
}
