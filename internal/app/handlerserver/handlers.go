// package handlerserver
//
// import (
//
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/Ippolid/shortLink/internal/app"
//	"github.com/Ippolid/shortLink/internal/models"
//	"github.com/gin-gonic/gin"
//	"io"
//	"net/http"
//	"strings"
//
// )
//
//	func (s *Server) PostCreate(c *gin.Context) {
//		val, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			c.String(http.StatusBadRequest, "Can't read body")
//			return
//		}
//
//		id, err := app.GenerateShortID(val)
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
//
// //	func (s *Server) GetID(res http.ResponseWriter, req *http.Request) {
// //		if req.Method != http.MethodGet {
// //			http.Error(res, "Not GET method", http.StatusBadRequest)
// //			return
// //		}
// //
// //		id := strings.TrimPrefix(req.URL.Path, "/")
// //		val, err := s.database.Data[id]
// //		if !err {
// //			http.Error(res, "Can`t find link", http.StatusBadRequest)
// //			return
// //		}
// //
// //		res.Header().Set("content-type", "text/plain")
// //		// устанавливаем код 200
// //
// //		http.Redirect(res, req, val, http.StatusTemporaryRedirect)
// //	}
//
//	func (s *Server) GetID(c *gin.Context) {
//		var val string
//		var err error
//		var exist bool
//		id := c.Param("id")
//		if s.Db == nil {
//			val, exist = s.database.Data[id]
//			if !exist {
//				c.String(http.StatusBadRequest, "Can't find link")
//				return
//			}
//		} else {
//			val, err = s.Db.GetLink(id)
//			if err != nil {
//				c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
//				return
//			}
//		}
//
//		fmt.Println(val)
//		if err != nil {
//			c.String(http.StatusBadRequest, "Can't find link")
//			return
//		}
//
//		c.Header("content-type", "text/plain")
//		c.Redirect(http.StatusTemporaryRedirect, val)
//	}
//
//	func (s *Server) PingDB(c *gin.Context) {
//		b, err := s.Db.Ping()
//		if err != nil {
//			c.String(http.StatusInternalServerError, "DB is not available")
//			return
//		}
//		if b {
//			c.Status(http.StatusOK)
//			return
//		}
//		c.Status(http.StatusInternalServerError)
//	}
//
// //	func ValidationMiddleware(next http.Handler) http.Handler {
// //		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// //			// Проверяем метод запроса
// //			if r.Method != http.MethodPost && r.Method != http.MethodGet {
// //				http.Error(w, "Method Not Allowed", http.StatusBadRequest)
// //				return
// //			}
// //
// //			// Если метод POST, проверяем, что путь "/"
// //			if r.Method == http.MethodPost && r.URL.Path != "/" {
// //				http.Error(w, "Invalid POST path", http.StatusBadRequest)
// //				return
// //			}
// //
// //			// Если метод GET, проверяем, что путь содержит ID
// //			if r.Method == http.MethodGet && (r.URL.Path == "/" || strings.Contains(r.URL.Path, "/ ")) {
// //				http.Error(w, "Invalid GET path", http.StatusBadRequest)
// //				return
// //			}
// //
// //			next.ServeHTTP(w, r)
// //		})
// //	}
//
//	func ValidationMiddleware() gin.HandlerFunc {
//		return func(c *gin.Context) {
//			// Проверяем метод запроса
//			if c.Request.Method != http.MethodPost && c.Request.Method != http.MethodGet {
//				c.String(http.StatusBadRequest, "Method Not Allowed")
//				c.Abort()
//				return
//			}
//
//			// Если метод POST, проверяем, что путь "/"
//			if c.Request.Method == http.MethodPost && c.Request.URL.Path != "/" {
//				c.String(http.StatusBadRequest, "Invalid POST path")
//				c.Abort()
//				return
//			}
//
//			// Если метод GET, проверяем, что путь содержит ID
//			if c.Request.Method == http.MethodGet && (c.Request.URL.Path == "/" || strings.Contains(c.Request.URL.Path, "/ ")) {
//				c.String(http.StatusBadRequest, "Invalid GET path")
//				c.Abort()
//				return
//			}
//
//			c.Next()
//		}
//	}
//
//	func (s *Server) PostAPI(c *gin.Context) {
//		var req models.PostRerquest
//		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
//			c.String(http.StatusBadRequest, "Invalid JSON data")
//			return
//		}
//
//		id := app.GenerateShortID([]byte(req.URL))
//		if s.Db == nil {
//			_, exist := s.database.Data[id]
//			if exist {
//				response := models.PostResponse{
//					Result: s.Adr + id,
//				}
//				c.JSON(http.StatusConflict, response)
//				return
//			}
//			s.database.SaveLink([]byte(req.URL), id)
//		} else {
//			err := s.Db.InsertLink(id, req.URL)
//			if err != nil {
//				fmt.Println(err)
//				if strings.Contains(err.Error(), "link exists") {
//					response := models.PostResponse{
//						Result: s.Adr + id,
//					}
//					c.JSON(http.StatusConflict, response)
//					return
//				}
//				c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
//				return
//			}
//		}
//		response := models.PostResponse{
//			Result: s.Adr + id,
//		}
//		c.Header("Content-Type", "application/json")
//		c.JSON(http.StatusCreated, response)
//	}
//
//	func (s *Server) PostBatch(c *gin.Context) {
//		var req []models.PostBatchReq
//		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
//			c.String(http.StatusBadRequest, "Invalid JSON data")
//			return
//		}
//		var otv models.PostBatchResp
//		var resp []models.PostBatchResp
//		for _, r := range req {
//			if r.ID != "" && r.URL != "" {
//				otv.ID = r.ID
//				k := app.GenerateShortID([]byte(r.URL))
//				otv.URL = s.Adr + k
//
//				if s.Db == nil {
//					s.database.SaveLink([]byte(r.URL), k)
//				} else {
//					err := s.Db.InsertLink(k, r.URL)
//					if err != nil {
//						c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
//					}
//
//				}
//				resp = append(resp, otv)
//			}
//		}
//
//		c.Header("Content-Type", "application/json")
//		c.JSON(http.StatusCreated, resp)
//
// }
//
//	func (s *Server) PostURL(c *gin.Context) {
//		// читаем запрос из body
//		body, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
//			return
//		}
//
//		// проверяем на пустой body
//		if len(body) == 0 {
//			c.JSON(http.StatusNotFound, gin.H{
//				"response": gin.H{
//					"text": "Извините, я пока ничего не умею",
//				},
//				"version": "1.0",
//			})
//			return
//		}
//
//		// получаем userID из контекста
//		userID, exists := c.Get("user_id")
//		if !exists || userID == "" {
//			c.String(http.StatusBadRequest, "Error = not userID")
//		}
//
//		// создаем короткую ссылку
//		encodeURL, err := s.database.SaveLink(string(body), userID.(string))
//		if err != nil {
//			c.String(http.StatusConflict, "text/plain; charset=utf-8", s.Adr+encodeURL)
//			return
//		}
//
//		// записываем заголовок, статус и короткую ссылку
//		c.String(http.StatusCreated, "text/plain; charset=utf-8", s.Adr+encodeURL)
//	}
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

	// Получение userID из контекста
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.String(http.StatusBadRequest, "Error = not userID")
		return
	}

	id := app.GenerateShortID(string(val), userID.(string))
	if s.Db == nil {
		_, exist := s.database.Data[id]
		if exist {
			c.String(http.StatusConflict, s.Adr+id)
			return
		}
		s.database.SaveLink(string(val), userID.(string))
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

func (s *Server) PostAPI(c *gin.Context) {
	var req models.PostRerquest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Получение userID из контекста
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.String(http.StatusBadRequest, "Error = not userID")
		return
	}

	id := app.GenerateShortID(req.URL, userID.(string))
	if s.Db == nil {
		_, exist := s.database.Data[id]
		if exist {
			response := models.PostResponse{
				Result: s.Adr + id,
			}
			c.JSON(http.StatusConflict, response)
			return
		}
		s.database.SaveLink(req.URL, userID.(string))
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

	// Получение userID из контекста
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.String(http.StatusBadRequest, "Error = not userID")
		return
	}

	var otv models.PostBatchResp
	var resp []models.PostBatchResp
	for _, r := range req {
		if r.ID != "" && r.URL != "" {
			otv.ID = r.ID
			k := app.GenerateShortID(r.URL, userID.(string))
			otv.URL = s.Adr + k

			if s.Db == nil {
				s.database.SaveLink(r.URL, userID.(string))
			} else {
				err := s.Db.InsertLink(k, r.URL)
				if err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
					return
				}
			}
			resp = append(resp, otv)
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, resp)
}

func (s *Server) UserUrls(c *gin.Context) {
	// Получение userId из контекста
	userID, exists := c.Get("user_id")
	if !exists || userID == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Создаем структуру для ответа
	var userLinks []models.GETUserLinks
	var resp models.GETUserLinks

	// Если в базе данных есть ссылки для пользователя
	if urls, exists := s.database.Users[userID.(string)]; exists && len(urls) > 0 {
		for _, url := range urls {
			// Находим id для каждого URL
			id := app.GenerateShortID(url, userID.(string))
			resp.ShortURL = s.Adr + id
			resp.OriginalURL = url
			userLinks = append(userLinks, resp)
		}
	}

	// Устанавливаем заголовок Content-Type для JSON
	c.Header("Content-Type", "application/json")

	// Возвращаем данные со статусом 200 даже если массив пустой
	c.JSON(http.StatusOK, userLinks)
}
