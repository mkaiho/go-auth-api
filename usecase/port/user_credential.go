package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type (
	UserCredentialCreateInput struct {
		UserID   entity.ID
		Email    entity.Email
		Password entity.HashedPassword
	}
	UserCredentialCreateUpdateInput struct {
		UserID   entity.ID
		Password entity.HashedPassword
	}
)

type UserCredentialGateway interface {
	GetByEmail(ctx context.Context, email entity.Email) (*entity.UserCredential, error)
	Check(ctx context.Context, email entity.Email, password entity.Password) error
	Create(ctx context.Context, input UserCredentialCreateInput) (*entity.UserCredential, error)
	Update(ctx context.Context, input UserCredentialCreateUpdateInput) (*entity.UserCredential, error)
}
