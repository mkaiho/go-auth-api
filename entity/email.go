package entity

import (
	"errors"
	"fmt"
)

type Email string

func ParseEmail(v string) (Email, error) {
	email := Email(v)
	if err := email.Validate(); err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}
	return Email(v), nil
}

func (e Email) Validate() error {
	if len(e) == 0 {
		return errors.New("empty")
	}
	return nil
}
