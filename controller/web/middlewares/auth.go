package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/controller/web/handlers"
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/mkaiho/go-auth-api/util"
)

func CheckAuth(txm port.TransactionManager, credGateway port.UserCredentialGateway) handlers.Handler {
	return func(gc *gin.Context) {
		var err error
		var auth *handlers.Auth
		ctx := gc.Request.Context()
		logger := util.GLogger()
		defer func() {
			if err != nil {
				logger.Error(err, "failed to check auth")
				gErr := gc.Error(err)
				if handlers.IsAuthError(err) {
					gErr.SetType(gin.ErrorTypePublic)
				}
				gc.Abort()
			}
		}()
		ctx, err = txm.BeginContext(ctx)
		defer txm.Rollback(ctx)
		if err != nil {
			return
		}
		auth, err = handlers.GetAuthInfo(gc)
		if err != nil {
			return
		}
		email, err := entity.ParseEmail(auth.User)
		if err != nil {
			return
		}
		password, err := entity.ParsePassword(auth.Password)
		if err != nil {
			return
		}
		err = credGateway.Check(ctx, email, password)
		if err != nil {
			return
		}
	}
}
