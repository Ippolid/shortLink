package app

import (
	"io"
	"net/http"
	"strings"
)

func (s *Server) PostCreate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Not POST method", http.StatusMethodNotAllowed)
		return
	}

	val, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Can`t read body", http.StatusUnprocessableEntity)
		return
	}

	id := GenerateShortID(val)
	s.database.SaveLink(val, id)

	res.Header().Set("content-type", "text/plain")
	// устанавливаем код 200
	res.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	_, err = res.Write([]byte(host + id))
	if err != nil {
		http.Error(res, "Can`t write body", http.StatusUnprocessableEntity)
		return
	}
}

func (s *Server) GetId(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Not GET method", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(req.URL.Path, "/")
	val, err := s.database.Data[id]
	if err != true {
		http.Error(res, "Can`t find link", http.StatusBadRequest)
		return
	}

	res.Header().Set("content-type", "text/plain")
	// устанавливаем код 200

	http.Redirect(res, req, val, http.StatusTemporaryRedirect)
}

func (s *Server) BadRequest(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "BadRequest", http.StatusBadRequest)
	return
}

func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusBadRequest)
			return
		}

		// Если метод POST, проверяем, что путь "/"
		if r.Method == http.MethodPost && r.URL.Path != "/" {
			http.Error(w, "Invalid POST path", http.StatusBadRequest)
			return
		}

		// Если метод GET, проверяем, что путь содержит ID
		if r.Method == http.MethodGet && (r.URL.Path == "/" || strings.Contains(r.URL.Path, "/ ")) {
			http.Error(w, "Invalid GET path", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
