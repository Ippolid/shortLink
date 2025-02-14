package app

import (
	"net/http"
)

type Server struct {
	database *Dbase
}

func New(st *Dbase) *Server {
	s := &Server{
		database: st,
	}
	return s
}

type MyHandler struct{}

var h MyHandler

func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data := []byte("Привет!")
	res.Write(data)
}

func (s *Server) newServer() *http.ServeMux {
	engine := http.NewServeMux()

	engine.Handle("/", ValidationMiddleware(http.HandlerFunc(s.PostCreate)))
	engine.Handle("/{id}", ValidationMiddleware(http.HandlerFunc(s.GetID)))
	//engine.HandleFunc("/{id}", s.GetId)

	return engine
}

func (s *Server) Start() error {
	err := http.ListenAndServe(`:8080`, s.newServer())
	if err != nil {
		return err
	}
	return nil
}
