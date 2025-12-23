package database

import (
	"errors"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations handles the up-migration logic
func RunMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})

	// Check if migrate driver is initialized
	if err != nil {
		return err
	}

	// Create new migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	// Check if migrations ran successfully
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
