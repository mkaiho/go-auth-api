package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/mkaiho/go-auth-api/entity"
	mocks "github.com/mkaiho/go-auth-api/mocks/usecase/port"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/stretchr/testify/mock"
)

var _ port.UserGateway = (*StubUserGateway)(nil)

type StubUserGateway struct {
	calledTimes int
	m           *mocks.UserGateway
	users       map[entity.ID]*entity.User
	createCall  *mock.Call
	mux         sync.RWMutex
}

func (g *StubUserGateway) Get(ctx context.Context, id entity.ID) (*entity.User, error) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	user, ok := g.users[id]
	if !ok {
		return nil, usecase.ErrNotFoundEntity
	}

	return user, nil
}

func (g *StubUserGateway) List(ctx context.Context, input port.UserListInput) (entity.Users, error) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	var users entity.Users
	for _, user := range g.users {
		if input.Email != nil && user.Email != *input.Email {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (g *StubUserGateway) Create(ctx context.Context, input port.UserCreateInput) (*entity.User, error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.calledTimes++
	id := entity.ID(fmt.Sprintf("%010d", g.calledTimes))
	g.createCall.
		Times(g.calledTimes).
		Return(&entity.User{
			ID:    id,
			Name:  input.Name,
			Email: input.Email,
		}, nil)

	user, err := g.m.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	g.users[user.ID] = user

	return user, nil
}

func (g *StubUserGateway) Update(ctx context.Context, input port.UserUpdateInput) (*entity.User, error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	for _, user := range g.users {
		if user.ID == input.ID {
			g.users[input.ID].Name = input.Name
			g.users[input.ID].Email = input.Email
			return g.users[input.ID], nil
		}
	}

	return nil, usecase.ErrNotFoundEntity
}

func NewStubUserGateway() port.UserGateway {
	userGateway := new(StubUserGateway)
	userGateway.m = new(mocks.UserGateway)
	userGateway.users = make(map[entity.ID]*entity.User)
	userGateway.createCall = userGateway.m.
		On("Create", mock.Anything, mock.Anything)

	return userGateway
}
