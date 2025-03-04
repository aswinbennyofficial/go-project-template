package main

import (
	"myapp/src/cassandra"
	"myapp/src/config"
	logs "myapp/src/log"
	"myapp/src/postgres"
	"myapp/src/rabbitmq"

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
		logger.Info().Msg("Running postgres migrations")
		if err := postgres.Migrate(pgConn, cfg.Postgres.Migrations.Path); err != nil {
			logger.Fatal().Err(err).Msg("Failed to run postgres migrations")
		}
		logger.Info().Msg("Database migrations completed successfully")
	}

    cassandraClient,err := cassandra.NewCassandraConnection(cfg.Cassandra,logger)
    if err!=nil{
        logger.Fatal().Err(err).Msg("Failed to connect to cassandra")
    }
    defer cassandraClient.Close()



    if cfg.Cassandra.Migrations.Enabled {
        err:=cassandra.Migrate(cassandraClient,logger, cfg.Cassandra.Migrations.Path, cfg.Cassandra.Keyspace)
        if err!=nil{
            logger.Fatal().Err(err).Msg("Failed to run cassandra migrations")
            // logger.Error().Err(err)
        }
    }


    redisClient,err := redis.NewRedisClientWrapper(cfg.Redis, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to Redis")
        
    }
    defer redisClient.Close()


    rabbitmqClient,err := rabbitmq.NewRabbitMQClientWrapper(cfg.RabbitMQ, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to RabbitMQ")
    }
    defer rabbitmqClient.Close()

    app := &config.App{
        Config: cfg,
        Logger: logger,
        Postgres: pgConn,
        Redis:  redisClient,
        Cassandra: cassandraClient,
        RabbitMQ: rabbitmqClient,
    }

    srv := server.NewServer(app)
    srv.Start()
}

