package app

import (
	"fmt"
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

	id := GenerateShortID(val)
	s.database.SaveLink(val, id)

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
	fmt.Println(id)
	val, exists := s.database.Data[id]
	if !exists {
		c.String(http.StatusBadRequest, "Can't find link")
		return
	}

	c.Header("content-type", "text/plain")
	c.Redirect(http.StatusTemporaryRedirect, val)
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
