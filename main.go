package main

import (
	"myapp/src/cassandra"
	"myapp/src/config"
	logs "myapp/src/log"
	"myapp/src/postgres"

	"myapp/src/redis"
	"myapp/src/server"
)

func main() {
    cfg,err := config.LoadConfig()
    if err != nil {
        panic(err)
    }
    
    logger := logs.NewLogger(cfg.Log)

    logger.Info().Msgf("Config %v", cfg)

    pgConn, err := postgres.NewPostgresConnection(cfg.Postgres, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer pgConn.Close()

    if cfg.Postgres.Migrations.Enabled {
		logger.Info().Msg("Running database migrations")
		if err := postgres.Migrate(pgConn, cfg.Postgres.Migrations.Path); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run database migrations")
		}
		logger.Info().Msg("Database migrations completed successfully")
	}

    cassandraClient,err := cassandra.NewCassandraConnection(cfg.Cassandra,logger)
    if err!=nil{
        logger.Fatal().Err(err).Msg("Failed to connect to cassandra")
    }
    defer cassandraClient.Close()


    redisClient,err := redis.NewRedisClient(cfg.Redis, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to Redis")
    }
    defer redisClient.Close()

    app := &config.App{
        Config: cfg,
        Logger: logger,
        Postgres: pgConn,
        Redis:  redisClient,
        Cassandra: cassandraClient,
    }

    srv := server.NewServer(app)
    srv.Start()
}

