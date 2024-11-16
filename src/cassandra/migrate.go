package cassandra

import (
	"fmt"
	"path/filepath"

	"github.com/gocql/gocql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"      
	"github.com/golang-migrate/migrate/v4/database/cassandra"
	"github.com/rs/zerolog"
)

// RunMigrations applies Cassandra migrations from a directory
func Migrate(session *gocql.Session, logger zerolog.Logger, migrationPath string, KeySpaceName string) error {
	// Convert migration path to absolute path
	absPath, err := filepath.Abs(migrationPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to resolve migration path")
		return fmt.Errorf("failed to resolve migration path: %w", err)
	}

	// Initialize Cassandra migration driver
	driver, err := cassandra.WithInstance(session, &cassandra.Config{
        KeyspaceName : KeySpaceName,
    })
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create Cassandra migration driver")
		return fmt.Errorf("failed to create Cassandra driver: %w", err)
	}

	// Initialize the migrate instance with file source
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", absPath), "cassandra", driver)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize migrations")
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Error().Err(err).Msg("Migration failed")
		return fmt.Errorf("migration failed: %w", err)
	}

	if err == migrate.ErrNoChange {
		logger.Info().Msg("No new Cassandra migrations to apply")
		return nil
	}

	logger.Info().Msg("Cassandra migrations applied successfully!")
	return nil
}
