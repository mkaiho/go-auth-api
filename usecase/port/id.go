package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type IDGenerator interface {
	Generate(ctx context.Context) (entity.ID, error)
}
