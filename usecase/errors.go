package usecase

import "errors"

var ErrNoAuthUser = errors.New("not exist auth user")
var ErrInvalidCredential = errors.New("invalid credential")

var ErrNotFoundEntity = errors.New("not found entity")
var ErrAlreadyExistsEntity = errors.New("already exists entity")
