package main

import (
	"myapp/src/config"
	"myapp/src/postgres"
	logs "myapp/src/log"

	"myapp/src/redis"
	"myapp/src/server"
)

func main() {
    cfg,err := config.LoadConfig()
    if err != nil {
        panic(err)
    }
    
    logger := logs.NewLogger(cfg.Log)

    dbConn, err := postgres.NewPostgresConnection(cfg.Postgres, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer dbConn.Close()

    if cfg.Postgres.Migrations.Enabled {
		logger.Info().Msg("Running database migrations")
		if err := postgres.Migrate(dbConn, cfg.Postgres.Migrations.Path); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run database migrations")
		}
		logger.Info().Msg("Database migrations completed successfully")
	}

    redisClient := redis.NewRedisClient(cfg.Redis, logger)
    defer redisClient.Close()

    app := &config.App{
        Config: cfg,
        Logger: logger,
        DB:     dbConn,
        Redis:  redisClient,
    }

    srv := server.NewServer(app)
    srv.Start()
}

