package web

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/util"
)

func NewGinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; len(raw) > 0 {
			path = path + "?" + raw
		}

		c.Next()

		logger := util.GLogger().
			WithValues("latency", fmt.Sprintf("%dµs", time.Since(start)/1000)).
			WithValues("clientIP", c.ClientIP()).
			WithValues("method", c.Request.Method).
			WithValues("statusCode", c.Writer.Status()).
			WithValues("path", path).
			WithValues("bodySize", c.Writer.Size()).
			WithValues()
		if msgs := c.Errors.ByType(gin.ErrorTypePrivate); len(msgs) > 0 {
			logger.Error(errors.New(msgs.String()), "request error")
			return
		}
		logger.Info("accepted request")
	}
}
