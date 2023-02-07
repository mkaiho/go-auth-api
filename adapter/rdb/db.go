package rdb

import (
	"context"
	"errors"

	"github.com/mkaiho/go-auth-api/util"
)

var ErrInvalidDriverName = errors.New("driver name is unknown")

type DriverName string

func (n DriverName) String() string {
	return string(n)
}

type DB interface {
	Begin() (Transaction, error)
}

type Config interface {
	GetDSN() string
	GetMaxConns() int
	GetDriverName() DriverName
}

func printQueryExecuted(ctx context.Context, query string, params ...interface{}) {
	logger := util.FromContext(ctx)
	logger.WithValues("query", query).WithValues("params", params).Debug("query executed")
}
