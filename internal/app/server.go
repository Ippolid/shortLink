package app

import (
	"github.com/gin-gonic/gin"
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

func (s *Server) newServer() *gin.Engine {
	engine := gin.New()
	engine.Use(ValidationMiddleware())

	engine.POST(
		"/",
		s.PostCreate,
	)
	engine.GET("/{id}",
		s.GetID,
	)

	engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "Route not found")
	})

	return engine
}

func (s *Server) Start() error {
	engine := s.newServer()
	return engine.Run(":8080")
}
