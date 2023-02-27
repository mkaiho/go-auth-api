package rdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkaiho/go-auth-api/entity"
	"github.com/mkaiho/go-auth-api/usecase/port"
)

var allUserColumns = []string{
	"id",
	"name",
	"email",
}

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

func (a *UserAccess) Get(ctx context.Context, tx Transaction, id entity.ID) (*UserRow, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM users WHERE id = ?",
		strings.Join(allUserColumns, ", "),
	)
	defer printQueryExecuted(ctx, query, id)

	var row UserRow
	err := tx.Get(ctx, &row, query, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (a *UserAccess) ListCount(ctx context.Context, tx Transaction, input port.UserListInput) (int, error) {
	query := "SELECT COUNT(id) FROM users"
	var where []string
	var args []interface{}
	if input.Email != nil {
		where = append(where, "email = ?")
		args = append(args, *input.Email)
	}
	if len(where) > 0 {
		query = query + " WHERE " + strings.Join(where, " AND")
	}
	defer printQueryExecuted(ctx, query, args...)

	var row int
	err := tx.Get(ctx, &row, query, args...)
	if err != nil {
		return row, err
	}

	return row, nil
}

func (a *UserAccess) List(ctx context.Context, tx Transaction, input port.UserListInput) ([]*UserRow, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM users",
		strings.Join(allUserColumns, ", "),
	)
	var where []string
	var args []interface{}
	if input.Email != nil {
		where = append(where, "email = ?")
		args = append(args, *input.Email)
	}
	if len(where) > 0 {
		query = query + " WHERE " + strings.Join(where, " AND")
	}
	defer printQueryExecuted(ctx, query, args...)

	var rows []*UserRow
	err := tx.Select(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
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

func (a *UserAccess) Update(ctx context.Context, tx Transaction, row *UserRow) error {
	query := "UPDATE users SET name = :name, email = :email WHERE id = :id"
	defer printQueryExecuted(ctx, query, row)

	_, err := tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}
