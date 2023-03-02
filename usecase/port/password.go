package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type PasswordManager interface {
	Hash(ctx context.Context, value string) (entity.HashedPassword, error)
	Compare(ctx context.Context, hashedPassword entity.HashedPassword, password entity.Password) error
}
