package adapter

import (
	"context"

	"github.com/mkaiho/go-auth-api/adapter/crypto"
	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase/port"
)

var _ (port.PasswordManager) = (*PasswordManager)(nil)

type PasswordManager struct {
	generator crypto.HashGenerator
}

func NewPasswordManager(generator crypto.HashGenerator) *PasswordManager {
	return &PasswordManager{
		generator: generator,
	}
}

func (pm *PasswordManager) Hash(ctx context.Context, value string) (entity.HashedPassword, error) {
	b, err := pm.generator.Generate(ctx, []byte(value))
	if err != nil {
		return "", err
	}
	pwd, err := entity.ParseHashedPassword(string(b))
	if err != nil {
		return "", err
	}
	return pwd, nil
}

func (pm *PasswordManager) Compare(ctx context.Context, hashedPassword entity.HashedPassword, password entity.Password) error {
	err := pm.generator.Compare(ctx, []byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
