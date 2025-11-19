package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) *migrator {
	return &migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *migrator) Up() error {
	err := goose.Up(m.db, m.migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
