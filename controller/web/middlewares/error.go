package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/usecase"
)

func NoMatchPathHandler() handlers.Handler {
	return func(gc *gin.Context) {
		gc.Error(usecase.ErrNotFoundEntity).SetType(gin.ErrorTypePublic)
	}
}
