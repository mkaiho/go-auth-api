package rdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkaiho/go-auth-api/entity"
)

var allUserCredentialColumns = []string{
	"c.id",
	"u.id user_id",
	"u.email",
	"c.password",
}

type UserCredentialRow struct {
	ID       string `db:"id" json:"id"`
	UserID   string `db:"user_id" json:"user_id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type UserCredentialAccess struct {
}

func NewUserCredential() *UserCredentialAccess {
	return &UserCredentialAccess{}
}

func (a *UserCredentialAccess) GetByUserID(ctx context.Context, tx Transaction, userID entity.ID) (*UserCredentialRow, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM user_credentials c INNER JOIN users u ON c.user_id = u.id WHERE u.id = ?",
		strings.Join(allUserCredentialColumns, ", "),
	)
	defer printQueryExecuted(ctx, query, userID)

	var row UserCredentialRow
	err := tx.Get(ctx, &row, query, userID)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (a *UserCredentialAccess) GetByEmail(ctx context.Context, tx Transaction, email entity.Email) (*UserCredentialRow, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM user_credentials c INNER JOIN users u ON c.user_id = u.id WHERE u.email = ?",
		strings.Join(allUserCredentialColumns, ", "),
	)
	defer printQueryExecuted(ctx, query, email)

	var row UserCredentialRow
	err := tx.Get(ctx, &row, query, email)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (a *UserCredentialAccess) Create(ctx context.Context, tx Transaction, row *UserCredentialRow) error {
	query := `
INSERT INTO user_credentials (id, user_id, password)
VALUES (:id, :user_id, :password)
`
	defer printQueryExecuted(ctx, query, row)

	_, err := tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}

func (a *UserCredentialAccess) UpdateByUserID(ctx context.Context, tx Transaction, row *UserCredentialRow) error {
	query := "UPDATE user_credentials SET password = :password WHERE user_id = :user_id"
	defer printQueryExecuted(ctx, query, UserCredentialRow{
		ID:       row.UserID,
		UserID:   row.UserID,
		Password: "*****",
	})

	_, err := tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}

func (a *UserCredentialAccess) Update(ctx context.Context, tx Transaction, row *UserCredentialRow) error {
	query := "UPDATE user_credentials SET password = :password WHERE user_id = :user_id"
	defer printQueryExecuted(ctx, query, row)

	_, err := tx.NamedExec(ctx, query, row)
	if err != nil {
		return err
	}

	return nil
}
