package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/controller/web/middlewares"
)

func NewGinServer() *gin.Engine {
	server := gin.New()
	server.Use(middlewares.NewGinLogger(), middlewares.Recovery())
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return server
}
