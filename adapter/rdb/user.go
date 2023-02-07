package rdb

import (
	"context"
)

type UserRow struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type UserAccess struct {
}

func NewUserAccess() *UserAccess {
	return &UserAccess{}
}

func (a *UserAccess) Create(ctx context.Context, tx Transaction, row *UserRow) error {
	query := `
INSERT INTO users (id, name, email)
VALUES (:id, :name, :email)
`
	defer printQueryExecuted(ctx, query, row)

	_, err := tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}
