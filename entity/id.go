package entity

import (
	"errors"
	"fmt"
)

type ID string

func ParseID(v string) (ID, error) {
	id := ID(v)
	if err := id.Validate(); err != nil {
		return "", fmt.Errorf("invalid id: %w", err)
	}
	return ID(v), nil
}

func (id ID) Validate() error {
	if len(id) == 0 {
		return errors.New("empty")
	}
	return nil
}

type IDs []ID
