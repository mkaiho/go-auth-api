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

var _ port.UserCredentialGateway = (*StubUserCredentialGateway)(nil)

type StubUserCredentialGateway struct {
	calledTimes int
	m           *mocks.UserCredentialGateway
	creds       map[entity.Email]*entity.UserCredential
	createCall  *mock.Call
	mux         sync.RWMutex
}

func (g *StubUserCredentialGateway) GetByEmail(ctx context.Context, email entity.Email) (*entity.UserCredential, error) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	for _, cred := range g.creds {
		if _, ok := g.creds[cred.Email]; ok {
			return cred, nil
		}
	}

	return nil, usecase.ErrNotFoundEntity

}

func (g *StubUserCredentialGateway) Create(ctx context.Context, input port.UserCredentialCreateInput) (*entity.UserCredential, error) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.calledTimes++
	id := entity.ID(fmt.Sprintf("%010d", g.calledTimes))
	g.createCall.
		Times(g.calledTimes).
		Return(&entity.UserCredential{
			ID:       id,
			UserID:   input.UserID,
			Email:    input.Email,
			Password: input.Password,
		}, nil)

	cred, err := g.m.Create(ctx, input)

	if err != nil {
		return nil, err
	}

	g.creds[cred.Email] = cred
	return cred, nil
}

func NewStubUserCredentialGateway() port.UserCredentialGateway {
	credGateway := new(StubUserCredentialGateway)
	credGateway.m = new(mocks.UserCredentialGateway)
	credGateway.creds = make(map[entity.Email]*entity.UserCredential)
	credGateway.createCall = credGateway.m.
		On("Create", mock.Anything, mock.Anything)
	return credGateway
}
