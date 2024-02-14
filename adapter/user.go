package adapter

import (
	"context"

	"github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/usecase/port"
)

var _ port.UserGateway = (*UserGateway)(nil)

type UserGateway struct {
	idgen      port.IDGenerator
	userAccess *rdb.UserAccess
}

func NewUserGateway(
	idgen port.IDGenerator,
	userAccess *rdb.UserAccess,
) *UserGateway {
	return &UserGateway{
		idgen:      idgen,
		userAccess: userAccess,
	}
}

func (g *UserGateway) Get(ctx context.Context, id entity.ID) (*entity.User, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	row, err := g.userAccess.Get(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	email, err := entity.ParseEmail(row.Email)
	if err != nil {
		return nil, err
	}
	user := entity.User{
		ID:    id,
		Name:  row.Name,
		Email: email,
	}

	return &user, nil
}

func (g *UserGateway) List(ctx context.Context, input port.UserListInput) (entity.Users, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := g.userAccess.List(ctx, tx, input)
	if err != nil {
		return nil, err
	}

	var users entity.Users
	for _, row := range rows {
		id, err := entity.ParseID(row.ID)
		if err != nil {
			return nil, err
		}
		email, err := entity.ParseEmail(row.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &entity.User{
			ID:    id,
			Name:  row.Name,
			Email: email,
		})
	}

	return users, nil
}

func (g *UserGateway) Create(ctx context.Context, input port.UserCreateInput) (*entity.User, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	count, err := g.userAccess.ListCount(ctx, tx, port.UserListInput{
		Email: &input.Email,
	})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, usecase.ErrAlreadyExistsEntity
	}

	id, err := g.idgen.Generate()
	if err != nil {
		return nil, err
	}
	created := entity.User{
		ID:    id,
		Name:  input.Name,
		Email: input.Email,
	}
	err = g.userAccess.Create(ctx, tx, &rdb.UserRow{
		ID:    created.ID.String(),
		Name:  created.Name,
		Email: created.Email.String(),
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (g *UserGateway) Update(ctx context.Context, input port.UserUpdateInput) (*entity.User, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = g.userAccess.Get(ctx, tx, input.ID)
	if err != nil {
		return nil, err
	}

	updated := entity.User{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}
	err = g.userAccess.Update(ctx, tx, &rdb.UserRow{
		ID:    updated.ID.String(),
		Name:  updated.Name,
		Email: updated.Email.String(),
	})
	if err != nil {
		return nil, err
	}

	return &updated, nil
}
