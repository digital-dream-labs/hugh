package sql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pressly/goose"
)

// New initializes a new database connection
func New(opts ...Option) (*sql.DB, error) {
	cfg := options{}

	for _, o := range opts {
		o(&cfg)
	}

	if cfg.errored() {
		return nil, fmt.Errorf("error during server setup: %v", cfg.errs)
	}

	var url string

	switch cfg.databaseType {
	case "mysql":
		url = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.username,
			cfg.password,
			cfg.host,
			cfg.port,
			cfg.name,
		)
	case "postgres":
		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%v/%s?sslmode=%s",
			cfg.username,
			cfg.password,
			cfg.host,
			cfg.port,
			cfg.name,
			cfg.tlsMode,
		)
	}

	db, err := sql.Open(cfg.databaseType, url)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %v", err)
	}

	return db, nil
}

// MigrateDatabase performs a database migration
func MigrateDatabase(opts ...Option) error {
	cfg := options{}

	for _, o := range opts {
		o(&cfg)
	}

	db, err := New(opts...)
	if err != nil {
		return err
	}

	// Mostly here in case the docker-compose db isn't up
	for loop := 0; loop < 5; loop++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close DB: %v\n", err)
		}
	}()

	if err := goose.SetDialect(cfg.databaseType); err != nil {
		return err
	}

	return goose.Up(db, fmt.Sprintf("database/%s", cfg.databaseType))
}
