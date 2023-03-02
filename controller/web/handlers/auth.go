package handlers

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/util"
)

var ErrNoAuthValue = errors.New("no auth value")
var ErrInvalidAuthValue = errors.New("invalid auth header value")
var ErrNotSupportedAuthType = errors.New("not supported auth type")

type Auth struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func GetAuthInfo(gc *gin.Context) (*Auth, error) {
	hValue := gc.Request.Header.Get("Authorization")
	if len(hValue) == 0 {
		return nil, ErrNoAuthValue
	}
	hValues := strings.SplitN(hValue, " ", 2)
	if len(hValues) != 2 {
		return nil, ErrInvalidAuthValue
	}

	authType := strings.TrimSpace(hValues[0])
	authValue := strings.TrimSpace(hValues[1])
	switch authType {
	default:
		return nil, ErrNotSupportedAuthType
	case "Basic":
		return getBasicAuthInfo(authValue)
	}
}

func IsAuthError(e error) bool {
	if errors.Is(e, ErrNoAuthValue) {
		return true
	}
	if errors.Is(e, ErrInvalidAuthValue) {
		return true
	}
	if errors.Is(e, ErrNotSupportedAuthType) {
		return true
	}
	if errors.Is(e, usecase.ErrNoAuthUser) {
		return true
	}
	if errors.Is(e, usecase.ErrInvalidCredential) {
		return true
	}
	return false
}

func getBasicAuthInfo(authValue string) (*Auth, error) {
	logger := util.GLogger()
	dec, err := base64.StdEncoding.DecodeString(authValue)
	if err != nil {
		logger.Error(err, "failed to decode auth value")
		return nil, ErrInvalidAuthValue
	}
	decValues := strings.Split(string(dec), ":")
	if len(decValues) != 2 {
		return nil, ErrInvalidAuthValue
	}

	auth := &Auth{
		User:     decValues[0],
		Password: decValues[1],
	}
	if len(auth.User) == 0 || len(auth.Password) == 0 {
		return nil, ErrInvalidAuthValue
	}

	return auth, nil
}
