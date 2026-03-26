package datastore

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func RunInTransaction(
	ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error,
) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Always defer rollback. If Commit is called, Rollback does nothing.
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func NewRunInTransaction[T any](
	ctx context.Context,
	db *sqlx.DB,
	fn func(*sqlx.Tx) (T, error),
) (T, error) {
	var result T
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return result, err
	}

	// Always defer rollback. If Commit is called, Rollback does nothing.
	defer tx.Rollback()

	result, err = fn(tx)
	if err != nil {
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, err
	}

	return result, nil
}
