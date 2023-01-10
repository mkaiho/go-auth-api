package port

import "github.com/mkaiho/go-auth-api/entity"

type UserCredentialRepository interface {
	GetByEmail(email entity.Email) (*entity.UserCredential, error)
	Create(email entity.Email, password entity.Password) error
}
