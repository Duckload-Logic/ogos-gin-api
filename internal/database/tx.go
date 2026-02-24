package database

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
