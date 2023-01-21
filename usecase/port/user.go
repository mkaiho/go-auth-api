package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type (
	UserCreateInput struct {
		Name  string
		Email entity.Email
	}
	UserListInput struct {
		Email *entity.Email
	}
)

type UserGateway interface {
	Get(ctx context.Context, id entity.ID) (*entity.User, error)
	List(ctx context.Context, input UserListInput) (entity.Users, error)
	Create(ctx context.Context, input UserCreateInput) (*entity.User, error)
}
