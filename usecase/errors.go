package usecase

import "errors"

var ErrNotFoundEntity = errors.New("not found entity")
var ErrAlreadyExistsEntity = errors.New("already exists entity")
