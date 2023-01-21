package interactor

import (
	"context"

	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/mkaiho/go-auth-api/util"
)

type (
	GetUserInput struct {
		ID entity.ID
	}
	FindUserInput struct {
		Email *entity.Email
	}
	CreateUserInput struct {
		Name     string
		Email    entity.Email
		Password entity.Password
	}
)

var _ UserInteractor = (*userInteractor)(nil)

type UserInteractor interface {
	GetUser(ctx context.Context, input GetUserInput) (*entity.User, error)
	FindUsers(ctx context.Context, input FindUserInput) (entity.Users, error)
	CreateUser(ctx context.Context, input CreateUserInput) (*entity.User, error)
}

type userInteractor struct {
	users     port.UserGateway
	userCreds port.UserCredentialGateway
}

func NewUserInteractor(
	users port.UserGateway,
	userCreds port.UserCredentialGateway,
) *userInteractor {
	return &userInteractor{
		users:     users,
		userCreds: userCreds,
	}
}

func (it *userInteractor) GetUser(
	ctx context.Context,
	input GetUserInput,
) (*entity.User, error) {
	logger := util.FromContext(ctx)

	user, err := it.users.Get(ctx, input.ID)
	if err != nil {
		logger.Error(err, "failed get user")
		return nil, err
	}

	return user, nil
}

func (it *userInteractor) FindUsers(
	ctx context.Context,
	input FindUserInput,
) (entity.Users, error) {
	logger := util.FromContext(ctx)

	users, err := it.users.List(ctx, port.UserListInput{
		Email: input.Email,
	})
	if err != nil {
		logger.Error(err, "failed find user")
		return nil, err
	}

	return users, nil
}

func (it *userInteractor) CreateUser(
	ctx context.Context,
	input CreateUserInput,
) (*entity.User, error) {
	logger := util.FromContext(ctx)

	user, err := it.users.Create(ctx, port.UserCreateInput{
		Name:  input.Name,
		Email: input.Email,
	})
	if err != nil {
		logger.Error(err, "failed create user")
		return nil, err
	}

	_, err = it.userCreds.Create(ctx, port.UserCredentialCreateInput{
		UserID:   user.ID,
		Email:    user.Email,
		Password: input.Password,
	})
	if err != nil {
		logger.Error(err, "failed create user credentials")
		return nil, err
	}

	return user, nil
}
