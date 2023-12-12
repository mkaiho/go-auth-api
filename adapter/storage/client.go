package storage

import (
	"context"
	"io"
)

type Client interface {
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	Save(ctx context.Context, path string, mime MimeType, body []byte) error
	Remove(ctx context.Context, path string) error
}
