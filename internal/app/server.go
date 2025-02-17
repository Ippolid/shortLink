package app

import (
	"github.com/gin-gonic/gin"
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

func (s *Server) newServer() *gin.Engine {
	engine := gin.New()

	engine.POST(
		"/",
		gin.WrapF(s.PostCreate),
	)
	engine.GET("/{id}", gin.WrapF(s.GetID))

	return engine
}

func (s *Server) Start() error {
	engine := s.newServer()
	return engine.Run(":8080")
}
