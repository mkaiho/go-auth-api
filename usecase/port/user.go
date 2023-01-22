package port

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
)

type (
	UserListInput struct {
		Email *entity.Email
	}
	UserCreateInput struct {
		Name  string
		Email entity.Email
	}
	UserUpdateInput struct {
		ID    entity.ID
		Name  string
		Email entity.Email
	}
)

type UserGateway interface {
	Get(ctx context.Context, id entity.ID) (*entity.User, error)
	List(ctx context.Context, input UserListInput) (entity.Users, error)
	Create(ctx context.Context, input UserCreateInput) (*entity.User, error)
	Update(ctx context.Context, input UserUpdateInput) (*entity.User, error)
}
