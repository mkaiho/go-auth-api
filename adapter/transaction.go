package adapter

import (
	"context"

	"github.com/mkaiho/go-auth-api/adapter/rdb"
	rdbAdapter "github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/usecase/port"
	"github.com/mkaiho/go-auth-api/util"
)

var _ port.TransactionManager = (*TransactionManager)(nil)

type TransactionManager struct {
	db *rdbAdapter.DB
}

func NewTransactionManager(rdb *rdbAdapter.DB) *TransactionManager {
	return &TransactionManager{
		db: rdb,
	}
}

func (tm *TransactionManager) BeginContext(parent context.Context) (context.Context, error) {
	ctx, err := rdbAdapter.ContextWithTx(parent, *tm.db)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (tm *TransactionManager) Rollback(ctx context.Context) error {
	rdbTx, err := rdbAdapter.TxFromContext(ctx)
	if err != nil {
		return err
	}
	if err := rdbTx.Rollback(); err != nil {
		return err
	}

	return nil
}

func (tm *TransactionManager) End(ctx context.Context) error {
	rdbTx, err := rdbAdapter.TxFromContext(ctx)
	if err != nil {
		return err
	}
	if err := rdbTx.Commit(); err != nil {
		return err
	}

	return nil
}

func DoInTx[T any](
	ctx context.Context, db rdb.DB,
	fn func(ctx context.Context, tx rdb.Transaction) (T, error),
) (result T, err error) {
	logger := util.FromContext(ctx)
	tx, err := db.Begin()
	if err != nil {
		return result, err
	}
	defer func(err *error) {
		if *err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				logger.Error(rbErr, "failed to rollback")
			}
			return
		}
		if cErr := tx.Commit(); cErr != nil {
			*err = cErr
		}
		logger.Debug("commited")
	}(&err)

	result, err = fn(ctx, tx)
	return result, err
}
