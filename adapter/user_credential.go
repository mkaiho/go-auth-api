package adapter

import (
	"context"
	"errors"

	"github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/mkaiho/go-auth-api/util"
)

var _ port.UserCredentialGateway = (*UserCredentialGateway)(nil)

type UserCredentialGateway struct {
	idgen           port.IDGenerator
	passwordManager port.PasswordManager
	userAccess      *rdb.UserAccess
	userCredAccess  *rdb.UserCredentialAccess
}

func NewUserCredentialGateway(
	idgen port.IDGenerator,
	passwordManager port.PasswordManager,
	userAccess *rdb.UserAccess,
	userCredAccess *rdb.UserCredentialAccess,
) port.UserCredentialGateway {
	return &UserCredentialGateway{
		idgen:           idgen,
		passwordManager: passwordManager,
		userAccess:      userAccess,
		userCredAccess:  userCredAccess,
	}
}

func (g *UserCredentialGateway) GetByEmail(ctx context.Context, email entity.Email) (*entity.UserCredential, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	row, err := g.userCredAccess.GetByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}

	id, err := entity.ParseID(row.ID)
	if err != nil {
		return nil, err
	}
	userID, err := entity.ParseID(row.UserID)
	if err != nil {
		return nil, err
	}
	pwd, err := entity.ParseHashedPassword(row.Password)
	if err != nil {
		return nil, err
	}
	userCred := entity.UserCredential{
		ID:       id,
		UserID:   userID,
		Email:    email,
		Password: pwd,
	}

	return &userCred, usecase.ErrNotFoundEntity

}

func (g *UserCredentialGateway) Create(ctx context.Context, input port.UserCredentialCreateInput) (*entity.UserCredential, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userRow, err := g.userAccess.Get(ctx, tx, input.UserID)
	if err != nil {
		return nil, err
	}
	_, err = g.userCredAccess.GetByUserID(ctx, tx, input.UserID)
	if err != nil {
		if !errors.Is(err, usecase.ErrNotFoundEntity) {
			return nil, err
		}
	}

	id, err := g.idgen.Generate(ctx)
	if err != nil {
		return nil, err
	}
	email, err := entity.ParseEmail(userRow.Email)
	if err != nil {
		return nil, err
	}
	created := entity.UserCredential{
		ID:       id,
		UserID:   input.UserID,
		Email:    email,
		Password: input.Password,
	}

	err = g.userCredAccess.Create(ctx, tx, &rdb.UserCredentialRow{
		ID:       created.ID.String(),
		UserID:   created.UserID.String(),
		Email:    created.Email.String(),
		Password: created.Password.String(),
	})
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (g *UserCredentialGateway) Update(ctx context.Context, input port.UserCredentialCreateUpdateInput) (*entity.UserCredential, error) {
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userRow, err := g.userAccess.Get(ctx, tx, input.UserID)
	if err != nil {
		return nil, err
	}
	credRow, err := g.userCredAccess.GetByUserID(ctx, tx, input.UserID)
	if err != nil {
		return nil, err
	}

	id, err := entity.ParseID(credRow.ID)
	if err != nil {
		return nil, err
	}
	email, err := entity.ParseEmail(userRow.Email)
	if err != nil {
		return nil, err
	}
	updated := entity.UserCredential{
		ID:       id,
		UserID:   input.UserID,
		Email:    email,
		Password: input.Password,
	}

	err = g.userCredAccess.Create(ctx, tx, &rdb.UserCredentialRow{
		ID:       updated.ID.String(),
		UserID:   updated.UserID.String(),
		Email:    updated.Email.String(),
		Password: updated.Password.String(),
	})
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (g *UserCredentialGateway) Check(ctx context.Context, email entity.Email, password entity.Password) error {
	logger := util.FromContext(ctx)
	tx, err := rdb.TxFromContext(ctx)
	if err != nil {
		return err
	}

	credRow, err := g.userCredAccess.GetByEmail(ctx, tx, email)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFoundEntity) {
			return usecase.ErrNoAuthUser
		}
		return err
	}
	hashed, err := entity.ParseHashedPassword(credRow.Password)
	if err != nil {
		return err
	}

	if err := g.passwordManager.Compare(ctx, hashed, password); err != nil {
		logger.Error(err, "failed to compare credentials")
		return usecase.ErrInvalidCredential
	}

	return nil
}
