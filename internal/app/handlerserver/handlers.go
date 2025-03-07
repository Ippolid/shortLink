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

//	func (s *Server) PostCreate(res http.ResponseWriter, req *http.Request) {
//		if req.Method != http.MethodPost {
//			http.Error(res, "Not POST method", http.StatusMethodNotAllowed)
//			return
//		}
//
//		val, err := io.ReadAll(req.Body)
//		if err != nil {
//			http.Error(res, "Can`t read body", http.StatusUnprocessableEntity)
//			return
//		}
//
//		id := GenerateShortID(val)
//		s.database.SaveLink(val, id)
//
//		res.Header().Set("content-type", "text/plain")
//		// устанавливаем код 200
//		res.WriteHeader(http.StatusCreated)
//		// пишем тело ответа
//		_, err = res.Write([]byte(host + id))
//		if err != nil {
//			http.Error(res, "Can`t write body", http.StatusUnprocessableEntity)
//			return
//		}
//	}
//
//	func (s *Server) PostCreate(c *gin.Context) {
//		val, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			c.String(http.StatusBadRequest, "Can't read body")
//			return
//		}
//
//		id := app.GenerateShortID(val)
//		if s.Db == nil {
//			_, exist := s.database.Data[id]
//			if exist {
//				c.String(http.StatusConflict, s.Adr+id)
//				return
//			}
//			s.database.SaveLink(val, id)
//		} else {
//			err = s.Db.InsertLink(id, string(val))
//			if err != nil {
//				fmt.Println(err)
//				if strings.Contains(err.Error(), "link exists") {
//					c.String(http.StatusConflict, s.Adr+id)
//					return
//				}
//				c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
//				return
//			}
//		}
//
//		c.Header("content-type", "text/plain")
//		c.String(http.StatusCreated, s.Adr+id)
//	}
func (s *Server) PostCreate(c *gin.Context) {
	// Получаем user_id из контекста (его устанавливает AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

	// Читаем тело запроса
	val, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't read body")
		return
	}

	// Генерируем уникальный short_id
	id := app.GenerateShortID(val)

	// Проверяем наличие базы данных
	if s.Db == nil {
		// Проверяем, есть ли уже такой short_id
		_, exist := s.database.Data[id]
		if exist {
			c.String(http.StatusConflict, s.Adr+id)
			return
		}
		s.database.SaveLink(val, id)
		// Сохраняем ссылку в локальную "базу"
		s.database.SaveUserLink(userIDStr, string(val))
	} else {
		// Сохраняем ссылку в БД (если она есть)
		err = s.Db.InsertLink(id, string(val), userIDStr)
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

	// Устанавливаем content-type и возвращаем результат
	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, s.Adr+id)
}

//	func (s *Server) GetID(res http.ResponseWriter, req *http.Request) {
//		if req.Method != http.MethodGet {
//			http.Error(res, "Not GET method", http.StatusBadRequest)
//			return
//		}
//
//		id := strings.TrimPrefix(req.URL.Path, "/")
//		val, err := s.database.Data[id]
//		if !err {
//			http.Error(res, "Can`t find link", http.StatusBadRequest)
//			return
//		}
//		res.Header().Set("content-type", "text/plain")
//
//		http.Redirect(res, req, val, http.StatusTemporaryRedirect)
//	}
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

//	func ValidationMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Проверяем метод запроса
//			if r.Method != http.MethodPost && r.Method != http.MethodGet {
//				http.Error(w, "Method Not Allowed", http.StatusBadRequest)
//				return
//			}
//
//			// Если метод POST, проверяем, что путь "/"
//			if r.Method == http.MethodPost && r.URL.Path != "/" {
//				http.Error(w, "Invalid POST path", http.StatusBadRequest)
//				return
//			}
//
//			// Если метод GET, проверяем, что путь содержит ID
//			if r.Method == http.MethodGet && (r.URL.Path == "/" || strings.Contains(r.URL.Path, "/ ")) {
//				http.Error(w, "Invalid GET path", http.StatusBadRequest)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем метод запроса
		if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
			c.String(http.StatusBadRequest, "Method Not Allowed")
			c.Abort()
			return
		}

		// Если метод POST, проверяем, что путь "/"
		if c.Request.Method == http.MethodPost && c.Request.URL.Path != "/" {
			c.String(http.StatusBadRequest, "Invalid POST path")
			c.Abort()
			return
		}

		// Если метод GET, проверяем, что путь содержит ID
		if c.Request.Method == http.MethodGet && (c.Request.URL.Path == "/" || strings.Contains(c.Request.URL.Path, "/ ")) {
			c.String(http.StatusBadRequest, "Invalid GET path")
			c.Abort()
			return
		}

		c.Next()
	}
}
func (s *Server) PostAPI(c *gin.Context) {
	var req models.PostRerquest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

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
		s.database.SaveUserLink(userIDStr, req.URL)
	} else {
		err := s.Db.InsertLink(id, req.URL, userIDStr)
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.String(http.StatusUnauthorized, "Unauthorized")
		return
	}
	userIDStr := userID.(string)

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
				s.database.SaveUserLink(userIDStr, r.URL)
			} else {
				err := s.Db.InsertLink(k, r.URL, userIDStr)
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
func (s *Server) GetUserURLs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var userURLs []string

	userIDStr := userID.(string)
	if s.Db == nil {
		userURLs, found := s.database.LoadUserLink(userIDStr)

		fmt.Println(userURLs)

		if !found || len(userURLs) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
	} else {
		userURLs, err := s.Db.GetLinksByUserID(userIDStr)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
		}

		fmt.Println(userURLs)

		if len(userURLs) == 0 {
			c.Status(http.StatusNoContent)
			return
		}

	}

	var otv models.UsersUrlResp
	var resp []models.UsersUrlResp
	var shortlink string

	for _, r := range userURLs {
		id := app.GenerateShortID([]byte(r))
		shortlink = s.Adr + id
		otv.ID = shortlink
		otv.URL = r
		resp = append(resp, otv)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, resp)
}
