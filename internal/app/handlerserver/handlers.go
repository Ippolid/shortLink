package handlerserver

import (
	"encoding/json"
	"fmt"
	"github.com/Ippolid/shortLink/config"
	"github.com/Ippolid/shortLink/internal/app"
	"github.com/Ippolid/shortLink/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID, exists := c.Get("userID")
	if !exists {
		c.String(http.StatusUnauthorized, "Can't get user ID")
		return
	}

	userIDStr, _ := userID.(string)

	if s.Db == nil {
		if _, exist := s.database.Data[id]; exist {
			c.String(http.StatusConflict, s.Adr+id)
			return
		}
		s.database.SaveLink(val, id)
		s.database.SaveUsersLink(userIDStr, id)
	} else {
		err = s.Db.InsertLink(id, string(val))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
			return
		}
	}

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, s.Adr+id)
}

//func (s *Server) PostCreate(c *gin.Context) {
//	val, err := io.ReadAll(c.Request.Body)
//	if err != nil {
//		c.String(http.StatusBadRequest, "Can't read body")
//		return
//	}
//
//	id := app.GenerateShortID(val)
//	userId, _ := app.GetUserId(c)
//	//if err != nil {
//	//	c.String(http.StatusUnauthorized, "Can't get user id")
//	//	return
//	//}
//	fmt.Println(userId)
//	if s.Db == nil {
//		_, exist := s.database.Data[id]
//		val2, _ := s.database.DataUsers[id]
//		if exist && val2 == userId {
//			c.String(http.StatusConflict, s.Adr+id)
//			return
//		}
//		s.database.SaveLink(val, id)
//		s.database.SaveUsersLink(userId, id)
//	} else {
//		err = s.Db.InsertLink(id, string(val))
//		if err != nil {
//			fmt.Println(err)
//			if strings.Contains(err.Error(), "link exists") {
//				c.String(http.StatusConflict, s.Adr+id)
//				return
//			}
//			c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
//			return
//		}
//	}
//
//	c.Header("content-type", "text/plain")
//	fmt.Println(s.database.DataUsers)
//	c.String(http.StatusCreated, s.Adr+id)
//}

//func (s *Server) GetID(c *gin.Context) {
//	var val string
//	var err error
//	var exist bool
//	id := c.Param("id")
//
//	if s.Db == nil {
//		val, exist = s.database.Data[id]
//		if !exist {
//			c.String(http.StatusBadRequest, "Can't find link")
//			return
//		}
//	} else {
//		val, err = s.Db.GetLink(id)
//		if err != nil {
//			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
//			return
//		}
//	}
//
//	fmt.Println(val)
//	if err != nil {
//		c.String(http.StatusBadRequest, "Can't find link")
//		return
//	}
//
//	c.Header("content-type", "text/plain")
//	c.Redirect(http.StatusTemporaryRedirect, val)
//}

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
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, userIDVal)
}

func (s *Server) PostAPI(c *gin.Context) {
	var req models.PostRerquest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	id := app.GenerateShortID([]byte(req.URL))
	userId, _ := app.GetUserId(c)
	//if err != nil {
	//	c.String(http.StatusUnauthorized, "Can't get user id")
	//	return
	//}
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
		s.database.SaveUsersLink(userId, id)
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
	userId, _ := app.GetUserId(c)
	//if err != nil {
	//	c.String(http.StatusUnauthorized, "Can't get user id")
	//	return
	//}
	for _, r := range req {
		if r.ID != "" && r.URL != "" {
			otv.ID = r.ID
			k := app.GenerateShortID([]byte(r.URL))
			otv.URL = s.Adr + k

			if s.Db == nil {
				s.database.SaveLink([]byte(r.URL), k)
				s.database.SaveUsersLink(userId, r.ID)
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

//	func (s *Server) UserUrls(c *gin.Context) {
//		var otv models.GETUserLinks
//		var resp []models.GETUserLinks
//		userId, _ := app.GetUserId(c)
//		//if err != nil {
//		//	c.String(http.StatusUnauthorized, "Can't get user id")
//		//	return
//		//}
//		for key, value := range s.database.DataUsers {
//			if value == userId {
//				otv.OriginalUrl = s.database.Data[key]
//				otv.ShortUrl = s.Adr + key
//				resp = append(resp, otv)
//			}
//		}
//		if len(resp) == 0 {
//			c.Header("Content-Type", "application/json")
//			c.JSON(http.StatusNoContent, gin.H{"message": "No links"})
//			return
//		} else {
//			c.Header("Content-Type", "application/json")
//			c.JSON(http.StatusOK, resp)
//		}
//	}
//func (s *Server) UserUrls(c *gin.Context) {
//	userIDVal, _ := c.Get("userID")
//	userID, _ := userIDVal.(string)
//
//	var resp []models.GETUserLinks
//
//	// Проходимся по вашему s.database.DataUsers,
//	// где key = "короткийID", value = "userID".
//	// Если value == userID, значит этот короткийID принадлежит данному пользователю.
//	for key, val := range s.database.DataUsers {
//		if val == userID {
//			// Заполняем структуру
//			var otv models.GETUserLinks
//			otv.OriginalUrl = s.database.Data[key] // допустим, тут исходный URL
//			otv.ShortUrl = s.Adr + key             // s.Adr = "http://localhost:8080/" (?)
//			resp = append(resp, otv)
//		}
//	}
//
//	if len(resp) == 0 {
//		// Нет ссылок
//		c.Header("Content-Type", "application/json")
//		c.JSON(http.StatusNoContent, gin.H{"message": "No links"})
//		return
//	}
//
//	// Если есть ссылки
//	c.Header("Content-Type", "application/json")
//	c.JSON(http.StatusOK, resp)
//}

//	func AuthMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			authHeader := c.GetHeader("Authorization")
//
//			// Будем искать префикс "Bearer "
//			var bearerToken string
//			if strings.HasPrefix(authHeader, "Bearer ") {
//				// Убираем префикс "Bearer "
//				bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
//			}
//
//			// Если токен пустой - возвращаем 401
//			if bearerToken == "" {
//				c.Status(http.StatusUnauthorized)
//				return
//			}
//
//			userID := app.GenerateShortID([]byte(bearerToken))
//			c.Set("user_id", userID)
//
//			// Иначе продолжаем - для примера просто вернем полученный токен в ответе
//
//			data, _ := json.Marshal(models.UserCookie{
//				UserID: userID,
//				Sign:   app.SignUserID(userID),
//			})
//
//			c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)
//
//			// Сохранение user_id в контексте
//			c.Set("userID", userID)
//			c.Next()
//		}
//	}
//func AuthMiddleware() gin.HandlerFunc {
//return func(c *gin.Context) {
//	cookieVal, err := c.Cookie(config.CookieName)
//	fmt.Println(cookieVal)
//	if err != nil {
//		// Куки нет -> генерируем нового userID
//		newUserID := app.GenerateShortID([]byte(uuid.New().String()))
//		sign := app.SignUserID(newUserID)
//
//		data, _ := json.Marshal(models.UserCookie{
//			UserID: newUserID,
//			Sign:   sign,
//		})
//		c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)
//
//		c.Set("userID", newUserID)
//	} else {
//		// Кука есть -> пытаемся распарсить JSON
//		var uc models.UserCookie
//		if err := json.Unmarshal([]byte(cookieVal), &uc); err != nil {
//			// Кука битая, генерируем заново
//			newUserID := app.GenerateShortID([]byte("fd"))
//			sign := app.SignUserID(newUserID)
//			data, _ := json.Marshal(models.UserCookie{UserID: newUserID, Sign: sign})
//
//			c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)
//			c.Set("userID", newUserID)
//		} else {
//			// Кука целая. Если вы пропускаете проверку подписи –
//			// просто считаем userID = uc.UserID, и этого достаточно:
//			c.Set("userID", uc.UserID)
//
//			// Если хотите удалить логику подписи, уберите SignUserID().
//			// Тогда храните в куке лишь userID, без sign.
//		}
//	}
//	c.Next()
//}

func (s *Server) GetID(c *gin.Context) {
	id := c.Param("id")

	// 1) Проверяем, не удалён ли
	if s.database.Deleted[id] {
		c.Status(http.StatusGone)
		return
	}

	// 2) Ищем URL
	var val string
	var exist bool
	var err error

	if s.Db == nil {
		val, exist = s.database.Data[id]
		if !exist {
			c.String(http.StatusBadRequest, "Can't find link")
			return
		}
	} else {
		val, err = s.Db.GetLink(id)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при обращении к дб: %v", err))
			return
		}
	}

	c.Redirect(http.StatusTemporaryRedirect, val)
}

// DeleteUserURLs - удаление коротких ссылок текущего пользователя
func (s *Server) DeleteUserURLs(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	userID, _ := userIDVal.(string)
	if userID == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	var shortIDs []string
	if err := json.NewDecoder(c.Request.Body).Decode(&shortIDs); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// "Асинхронная" пометка удалёнными
	go func() {
		for _, sid := range shortIDs {
			if s.database.DataUsers[sid] == userID {
				s.database.Deleted[sid] = true
			}
		}
	}()

	c.Status(http.StatusAccepted)
}

// UserUrls - GET /api/user/urls - возвращает все ссылки пользователя
func (s *Server) UserUrls(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.Status(http.StatusUnauthorized)
		return
	}

	userID, _ := userIDVal.(string)
	if userID == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	var resp []models.GETUserLinks

	for shortID, uid := range s.database.DataUsers {
		if uid == userID {
			resp = append(resp, models.GETUserLinks{
				ShortUrl:    s.Adr + shortID,
				OriginalUrl: s.database.Data[shortID],
			})
		}
	}

	c.Header("Content-Type", "application/json")

	if len(resp) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No links"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Пример AuthMiddleware, который даёт куку, если её нет или битая.
// Если нет Bearer-токена, userID = "" => пользователь гость
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var bearerToken = ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			bearerToken = strings.TrimSpace(authHeader[len("Bearer "):])
		}

		cookieVal, err := c.Cookie(config.CookieName)

		if err != nil {
			var sign string
			var newUserID string
			// Куки нет -> создаём новую
			if bearerToken != "" {
				newUserID = app.GenerateShortID([]byte(bearerToken))
				sign = app.SignUserID(newUserID)
			} else {
				newUserID := app.GenerateShortID([]byte(uuid.New().String()))
				sign = app.SignUserID(newUserID)
			}

			data, _ := json.Marshal(models.UserCookie{
				UserID: newUserID,
				Sign:   sign,
			})

			c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)

			c.Set("userID", newUserID)
		} else {
			// Кука есть -> разбираем её
			var uc models.UserCookie
			if err := json.Unmarshal([]byte(cookieVal), &uc); err != nil {
				// Битая кука -> пересоздаём
				newUserID := app.GenerateShortID([]byte(uuid.New().String()))
				sign := app.SignUserID(newUserID)

				data, _ := json.Marshal(models.UserCookie{
					UserID: newUserID,
					Sign:   sign,
				})

				c.SetCookie(config.CookieName, string(data), 3600*24, "/", "", false, true)
				c.Set("userID", newUserID)
			} else {
				// Кука нормальная -> читаем userID
				c.Set("userID", uc.UserID)
			}
		}
		fmt.Println(c.Get("userID"))
		c.Next()
	}
}
