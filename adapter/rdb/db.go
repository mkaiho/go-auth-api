package rdb

import (
	"errors"
	// _ "github.com/go-sql-driver/mysql"
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
