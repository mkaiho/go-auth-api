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

func (p Password) Validate() error {
	if len(p) == 0 {
		return errors.New("empty")
	}
	return nil
}

type UserCredential struct {
	ID       ID
	UserID   ID
	Email    Email
	Password Password
}
