package port

import "context"

type TransactionManager interface {
	BeginContext(parent context.Context) (context.Context, error)
	End(ctx context.Context) error
	Rollback(ctx context.Context) error
}
