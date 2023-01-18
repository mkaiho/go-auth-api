package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type (
	UserCredentialCreateInput struct {
		UserID   entity.ID
		Email    entity.Email
		Password entity.Password
	}
)

type UserCredentialGateway interface {
	GetByEmail(ctx context.Context, email entity.Email) (*entity.UserCredential, error)
	Create(ctx context.Context, input UserCredentialCreateInput) error
}
