package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinServer() *gin.Engine {
	server := gin.New()
	server.Use(NewGinLogger(), Recovery())
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return server
}
