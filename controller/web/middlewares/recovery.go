package middlewares

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/util"
)

func Recovery() handlers.Handler {
	logger := util.GLogger().WithCallDepth(2)
	return func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				var brokenPipe bool
				if ne, ok := p.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				headersToStr := strings.Join(headers, "\r\n")
				var pErr error
				if _pErr, ok := p.(error); ok {
					pErr = _pErr
				} else {
					pErr = fmt.Errorf(fmt.Sprintf("%v", _pErr))
				}
				if brokenPipe {
					logger.
						WithValues("headers", headersToStr).
						Error(pErr, "broken pipe")
				} else if gin.IsDebugging() {
					logger.
						WithValues("headers", headersToStr).
						Error(pErr, "panic recovered")
				} else {
					logger.Error(pErr, "panic recovered")
				}
				if brokenPipe {
					c.Error(pErr)
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
			if errMsgs := c.Errors.ByType(gin.ErrorTypeBind); len(errMsgs) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": errMsgs[0].Err.Error(),
				})
			} else if errMsgs := c.Errors.ByType(gin.ErrorTypePublic); len(errMsgs) > 0 {
				var msg string
				code := http.StatusBadRequest
				if errors.Is(errMsgs[0].Err, usecase.ErrNotFoundEntity) {
					code = http.StatusNotFound
				} else if errors.Is(errMsgs[0].Err, usecase.ErrAlreadyExistsEntity) {
					code = http.StatusConflict
				} else if handlers.IsAuthError(errMsgs[0].Err) {
					code = http.StatusUnauthorized
					msg = errMsgs[0].Err.Error()
				}

				if len(msg) == 0 {
					msg = http.StatusText(code)
				}
				c.AbortWithStatusJSON(code, gin.H{
					"message": msg,
				})
			}
		}()
		c.Next()
	}
}
