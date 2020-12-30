package sql

import (
	"database/sql"

	"github.com/pressly/goose"
)

// RunMySQLMigration is only used for tests
func RunMySQLMigration(db *sql.DB, dialect, path string) error {
	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	return goose.Up(db, path)
}
