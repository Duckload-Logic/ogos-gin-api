package datastore

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DB is a common interface for sqlx.DB and sqlx.Tx to provide
// unified database operations.
type DB interface {
	sqlx.ExtContext
	SelectContext(
		ctx context.Context,
		dest interface{},
		query string,
		args ...interface{},
	) error
	GetContext(
		ctx context.Context,
		dest interface{},
		query string,
		args ...interface{},
	) error
	NamedExecContext(
		ctx context.Context,
		query string,
		arg interface{},
	) (sql.Result, error)
}

func RunInTransaction(
	ctx context.Context, db *sqlx.DB, fn func(DB) error,
) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Always defer rollback. If Commit is called, Rollback does nothing.
	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func NewRunInTransaction[T any](
	ctx context.Context,
	db *sqlx.DB,
	fn func(DB) (T, error),
) (T, error) {
	var result T
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return result, err
	}

	// Always defer rollback. If Commit is called, Rollback does nothing.
	defer func() {
		_ = tx.Rollback()
	}()

	result, err = fn(tx)
	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}
