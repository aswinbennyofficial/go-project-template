package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// Migrate performs database migrations using goose
func Migrate(dbConn *pgxpool.Pool, migrationPath string) error {
	goose.SetBaseFS(nil)

    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }

    db := stdlib.OpenDBFromPool(dbConn)
    if err := goose.Up(db, migrationPath); err != nil {
        return fmt.Errorf("failed to apply migrations: %w", err)
    }
    if err := db.Close(); err != nil {
        return fmt.Errorf("failed to close database connection: %w", err)
    }

    return nil
}