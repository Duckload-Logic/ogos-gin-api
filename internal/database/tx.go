package database

import (
	"context"
	"database/sql"
)

func RunInTransaction(
	ctx context.Context, db *sql.DB, fn func(*sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, nil)
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
