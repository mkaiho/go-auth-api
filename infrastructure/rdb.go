package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/mkaiho/go-auth-api/adapter/rdb"
	"github.com/mkaiho/go-auth-api/usecase"
)

const (
	DriverNameUnknown rdb.DriverName = ""
	DriverNameMySQL   rdb.DriverName = "mysql"
)

var driverNames = []rdb.DriverName{
	DriverNameMySQL,
}

func (c *MySQLConfig) GetDriverName() rdb.DriverName {
	return DriverNameMySQL
}

// RDB
var _ rdb.DB = (*RDB)(nil)

type RDB struct {
	db *sqlx.DB
}

func (db *RDB) Begin() (rdb.Transaction, error) {
	tx, err := db.db.Beginx()
	if err != nil {
		return nil, err
	}

	return &RDBTransaction{
		tx: tx,
	}, err
}

func OpenRDB(conf rdb.Config) (*RDB, error) {
	driverName := conf.GetDriverName()
	if driverName == DriverNameUnknown {
		return nil, rdb.ErrInvalidDriverName
	}

	db, err := sqlx.Open(driverName.String(), conf.GetDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conf.GetMaxConns())

	return &RDB{
		db: db,
	}, nil
}

// RDBTransaction
var _ (rdb.Transaction) = (*RDBTransaction)(nil)

type RDBTransaction struct {
	tx *sqlx.Tx
}

func (rt *RDBTransaction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if err := rt.tx.GetContext(ctx, dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usecase.ErrNotFoundEntity
		}
		return err
	}

	return nil
}

func (rt *RDBTransaction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return rt.tx.SelectContext(ctx, dest, query, args...)
}

func (rt *RDBTransaction) NamedExec(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return rt.tx.NamedExecContext(ctx, query, arg)
}

func (rt *RDBTransaction) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return rt.tx.ExecContext(ctx, query, args...)
}

func (rt *RDBTransaction) Commit() error {
	return rt.tx.Commit()
}

func (rt *RDBTransaction) Rollback() error {
	return rt.tx.Rollback()
}
