package postgres

import (
	"context"
	"fmt"
	"myapp/src/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)


func NewPostgresConnection(config config.PostgresConfig, log zerolog.Logger) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=10",
        config.Host, config.Port, config.User, config.Password, config.DBName)
    
    log.Info().Msgf("Attempting to connect to database at %s:%d", config.Host, config.Port)
    
    // Initial delay to allow database to start up
    time.Sleep(30 * time.Second)
    
    var db *pgxpool.Pool
    var err error
    
    // Retry logic for database connection
    for i := 0; i < 60; i++ { // Retry up to 60 times (10 minutes total)
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        db, err = pgxpool.New(ctx, dsn)
        cancel()
        if err == nil {
            break
        }
        log.Error().Err(err).Msgf("Failed to connect to database, retrying (%d/60)...", i+1)
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }
    
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database after retries")
        return nil, err
    }
    
    // Test the connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    err = db.Ping(ctx)
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to ping database")
        return nil, err
    }
    
    log.Info().Msg("Connected to database successfully")
    return db, nil
}