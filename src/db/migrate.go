package db

import (
    "database/sql"
    "github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
    goose.SetBaseFS(nil)
    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }
    return goose.Up(db, "src/db/migrations")
}