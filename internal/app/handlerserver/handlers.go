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
func (s *Server) PostCreate(c *gin.Context) {
	val, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Can't read body")
		return
	}

	id := app.GenerateShortID(val)
	//s.database.SaveLink(val, id)
	//fmt.Println(s.database)
	err = s.Db.InsertLink(id, string(val))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Can't save link: %v", err))
		return
	}

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
//
//		res.Header().Set("content-type", "text/plain")
//		// устанавливаем код 200
//
//		http.Redirect(res, req, val, http.StatusTemporaryRedirect)
//	}
func (s *Server) GetID(c *gin.Context) {
	id := c.Param("id")
	//val, exists := s.database.Data[id]
	val, err := s.Db.GetLink(id)
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

	id := app.GenerateShortID([]byte(req.URL))
	//s.database.SaveLink([]byte(req.URL), id)
	err := s.Db.InsertLink(id, req.URL)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Ошибка при вставке данных в дб: %v", err))
		return
	}

	response := models.PostResponse{
		Result: s.Adr + id,
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, response)
}
