package port

import (
	"github.com/mkaiho/go-auth-api/entity"
)

type IDGenerator interface {
	Generate() (entity.ID, error)
}
