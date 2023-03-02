package entity

import (
	"errors"
	"fmt"
)

type Password string

func ParsePassword(v string) (Password, error) {
	id := Password(v)
	if err := id.Validate(); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return Password(v), nil
}

func (p Password) String() string {
	return string(p)
}

func (p Password) Validate() error {
	if len(p) == 0 {
		return errors.New("empty")
	}
	return nil
}

type HashedPassword string

func ParseHashedPassword(v string) (HashedPassword, error) {
	id := HashedPassword(v)
	if err := id.Validate(); err != nil {
		return "", fmt.Errorf("invalid hashed password: %w", err)
	}
	return HashedPassword(v), nil
}

func (p HashedPassword) String() string {
	return string(p)
}

func (p HashedPassword) Validate() error {
	if len(p) == 0 {
		return errors.New("empty")
	}
	return nil
}

type UserCredential struct {
	ID       ID
	UserID   ID
	Email    Email
	Password HashedPassword
}

type UserCredentials []*UserCredential
