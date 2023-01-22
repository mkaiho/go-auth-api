package web

import (
	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/controller/web/middlewares"
	"github.com/mkaiho/go-auth-api/controller/web/routes"
)

type Server struct {
	e *gin.Engine
}

func (s *Server) Use(middleware ...handlers.Handler) {
	for _, m := range middleware {
		s.e.Use(gin.HandlerFunc(m))
	}
}

func (s *Server) Handle(httpMethod string, relativePath string, h ...handlers.Handler) {
	var hs []gin.HandlerFunc
	for _, handler := range h {
		hs = append(hs, gin.HandlerFunc(handler))
	}
	s.e.Handle(httpMethod, relativePath, hs...)
}

func (s *Server) Run(addr ...string) error {
	return s.e.Run(addr...)
}

func NewGinServer(r ...*routes.Route) *Server {
	server := &Server{
		e: gin.New(),
	}
	server.Use(middlewares.NewGinLogger(), middlewares.Recovery())
	for _, route := range r {
		server.Handle(route.Method(), route.Path(), route.Handlers()...)
	}
	server.Use(middlewares.NoMatchPathHandler())

	return server
}
