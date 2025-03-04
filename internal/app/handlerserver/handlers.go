package handlerserver

import (
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func (s *Server) PostCreate(c *gin.Context) {
	val, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't read body")
		return
	}

	id := app.GenerateShortID(val)
	if s.Db == nil {
		_, exist := s.database.Data[id]
		if exist {
			c.String(http.StatusConflict, s.Adr+id)
			return
		}
		s.database.SaveLink(val, id)
	} else {
		err = s.Db.InsertLink(id, string(val))
		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "link exists") {
				c.String(http.StatusConflict, s.Adr+id)
				return
			}
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
			return
		}
	}

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, s.Adr+id)
}

func (s *Server) GetID(c *gin.Context) {
	var val string
	var err error
	var exist bool
	id := c.Param("id")
	if s.Db == nil {
		val, exist = s.database.Data[id]
		if !exist {
			c.String(http.StatusBadRequest, "Can't find link")
			return
		}
	} else {
		val, err = s.Db.GetLink(id)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
			return
		}
	}

	fmt.Println(val)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't find link")
		return
	}

	c.Header("content-type", "text/plain")
	c.Redirect(http.StatusTemporaryRedirect, val)
}

func (s *Server) PingDB(c *gin.Context) {
	b, err := s.Db.Ping()
	if err != nil {
		c.String(http.StatusInternalServerError, "DB is not available")
		return
	}
	if b {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusInternalServerError)
}

func (s *Server) TestCookie(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	// Будем искать префикс "Bearer "
	var bearerToken string
	if strings.HasPrefix(authHeader, "Bearer ") {
		// Убираем префикс "Bearer "
		bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// Если токен пустой - возвращаем 401
	if bearerToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no valid bearer token found",
		})
		return
	}

	// Иначе продолжаем - для примера просто вернем полученный токен в ответе
	c.JSON(http.StatusOK, gin.H{
		"bearer_token": bearerToken,
	})
}

func (s *Server) PostAPI(c *gin.Context) {
	var req models.PostRerquest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	id := app.GenerateShortID([]byte(req.URL))
	if s.Db == nil {
		_, exist := s.database.Data[id]
		if exist {
			response := models.PostResponse{
				Result: s.Adr + id,
			}
			c.JSON(http.StatusConflict, response)
			return
		}
		s.database.SaveLink([]byte(req.URL), id)
	} else {
		err := s.Db.InsertLink(id, req.URL)
		if err != nil {
			fmt.Println(err)
			if strings.Contains(err.Error(), "link exists") {
				response := models.PostResponse{
					Result: s.Adr + id,
				}
				c.JSON(http.StatusConflict, response)
				return
			}
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
			return
		}
	}
	response := models.PostResponse{
		Result: s.Adr + id,
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}

func (s *Server) PostBatch(c *gin.Context) {
	var req []models.PostBatchReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}
	var otv models.PostBatchResp
	var resp []models.PostBatchResp
	for _, r := range req {
		if r.ID != "" && r.URL != "" {
			otv.ID = r.ID
			k := app.GenerateShortID([]byte(r.URL))
			otv.URL = s.Adr + k

			if s.Db == nil {
				s.database.SaveLink([]byte(r.URL), k)
			} else {
				err := s.Db.InsertLink(k, r.URL)
				if err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
				}

			}
			resp = append(resp, otv)
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, resp)
}

//	func AuthMiddleware() gin.HandlerFunc {
//		//return func(c *gin.Context) {
//		//	cookie, err := c.Cookie(config.CookieName)
//		//	var userID string
//		//
//		//	if err != nil || userID, _ = app.ValidateCookie(cookie); err != nil {
//		//		// Генерация нового user_id
//		//		userID = uuid.NewString()
//		//		db.Create(&User{ID: userID})
//		//
//		//		// Создание подписанной куки
//		//		data, _ := json.Marshal(map[string]string{
//		//			"user_id": userID,
//		//			"sign":    signUserID(userID),
//		//		})
//		//
//		//		c.SetCookie(cookieName, string(data), 3600*24, "/", "", false, true)
//		//	}
//		//
//		//	// Сохранение user_id в контексте
//		//	c.Set("user_id", userID)
//		//	c.Next()
//		return func(c *gin.Context) {
//			cookie, err := c.Cookie(config.CookieName)
//			var userID string
//
//			// Если кука существует, проверяем её
//			if err == nil {
//				userID, _ = app.ValidateCookie(cookie)
//			}
//
//			// Если куки нет или она невалидна, создаем новую
//			if err != nil || userID == "" {
//				userID = uuid.NewString()
//
//				// Создание подписанной куки
//				data, _ := json.Marshal(models.UserCookie{
//					UserID: userID,
//					Sign:   app.SignUserID(userID),
//				})
//
//				c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)
//			}
//
//			// Сохранение user_id в контексте
//			c.Set("user_id", userID)
//			c.Next()
//		}
//	}
