package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ShouldBind[T any](gc *gin.Context, obj T) error {
	vErr := new(validator.ValidationErrors)
	if err := gc.ShouldBindUri(obj); err != nil && !errors.As(err, vErr) {
		return err
	}
	if err := gc.ShouldBind(obj); err != nil && !errors.As(err, vErr) {
		return err
	}
	if binding.Validator != nil {
		return binding.Validator.ValidateStruct(obj)
	}

	return nil
}
