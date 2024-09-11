package db

import (
    "context"
    "fmt"
    "myapp/src/utils"
    "time"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/rs/zerolog"
)

func NewPostgresConnection(config utils.PostgresConfig, log zerolog.Logger) (*pgxpool.Pool, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=10", 
        config.Host, config.Port, config.User, config.Password, config.DBName)

    var db *pgxpool.Pool
    var err error

    // Retry logic for database connection
    for i := 0; i < 30; i++ { // Retry up to 30 times (5 minutes total)
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        db, err = pgxpool.Connect(ctx, dsn)
        cancel()

        if err == nil {
            break
        }

        log.Warn().Err(err).Msgf("Failed to connect to database, retrying (%d/30)...", i+1)
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