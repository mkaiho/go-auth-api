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
	"github.com/mkaiho/go-auth-api/util"
)

func Recovery() gin.HandlerFunc {
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
		}()
		c.Next()
	}
}
