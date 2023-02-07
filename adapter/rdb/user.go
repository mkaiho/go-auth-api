package rdb

import "context"

type UserRow struct {
	ID    string `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
}

type UserAccess struct {
	tx Transaction
}

func NewUserAccess(tx Transaction) *UserAccess {
	return &UserAccess{
		tx: tx,
	}
}

func (a *UserAccess) Create(ctx context.Context, row *UserRow) error {
	query := `
INSERT INTO users (id, name, email)
VALUES (:id, :name, :email)
`
	_, err := a.tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}
