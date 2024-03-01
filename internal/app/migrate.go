package app

import (
	"context"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

const migrationsPath = "migrations"

//go:embed migrations/*.sql
var fs embed.FS

// startMigrate executes database migrations.
func (a *App) startMigrate(ctx context.Context, migratePath string, dbName string, db *sqlx.DB) error {
	// Check if the database connection is alive
	err := db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("db connection not alive: %w", err)
	}

	// Create the migration database driver
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{
		DatabaseName: dbName,
		SchemaName:   "public",
	})
	if err != nil {
		return fmt.Errorf("db migration database driver error: %w", err)
	}

	// Create the migration source driver
	source, err := iofs.New(fs, migratePath)
	if err != nil {
		return fmt.Errorf("db migration source driver error: %w", err)
	}

	// Create a new migration instance
	instance, err := migrate.NewWithInstance("fs", source, dbName, driver)
	if err != nil {
		return fmt.Errorf("db migration instance error: %w", err)
	}

	// Execute the migrations
	if err := instance.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("db migration up error: %w", err)
	}

	return nil
}
