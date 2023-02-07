package adapter

import (
	"context"

	"github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/util"
)

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
