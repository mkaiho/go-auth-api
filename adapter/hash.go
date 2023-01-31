package adapter

import "context"

type HashGenerator interface {
	Generate(ctx context.Context, value []byte) ([]byte, error)
	Compare(ctx context.Context, value []byte, hashed []byte) error
}
