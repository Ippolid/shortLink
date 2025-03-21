package handlerserver

import (
	"github.com/Ippolid/shortLink/internal/app/middleware"
	"github.com/Ippolid/shortLink/internal/app/storage"
	"github.com/Ippolid/shortLink/internal/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	database *storage.Dbase
	Host     string
	Adr      string
	Db       *storage.DataBase
}

func New(st *storage.Dbase, adr, host string, db *storage.DataBase) *Server {
	s := &Server{
		database: st,
		Host:     host,
		Adr:      adr,
		Db:       db,
	}
	return s
}

func (s *Server) newServer() *gin.Engine {
	engine := gin.New()
	//engine.Use(ValidationMiddleware())
	engine.Use(logger.RequestLogger())
	engine.Use(gzipDecompressMiddleware()) // Декомпрессия входящих запросов
	engine.Use(gzipMiddleware())

	engine.POST(
		"/",
		middleware.AuthMiddleware(),
		s.PostCreate,
	)
	engine.GET("/:id",
		s.GetID,
	)
	engine.POST("/api/shorten",
		middleware.AuthMiddleware(),
		s.PostAPI,
	)

	engine.GET("/ping",
		s.PingDB,
	)

	engine.POST("/api/shorten/batch",
		middleware.AuthMiddleware(),
		s.PostBatch,
	)

	engine.GET("api/user/urls",
		middleware.AuthMiddleware(),
		s.GetUserURLs,
	)

	engine.DELETE("/api/user/urls",
		middleware.AuthMiddleware(),
		s.DeleteLinks)

	engine.NoRoute(func(c *gin.Context) {
		c.String(http.StatusBadRequest, "Route not found")
	})

	return engine
}

func (s *Server) Start() error {
	engine := s.newServer()
	return engine.Run(s.Host)
}
