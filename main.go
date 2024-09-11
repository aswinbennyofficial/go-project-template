package main

import (
	"log"
	"myapp/src/db"
	logs "myapp/src/log"

	"myapp/src/redis"
	"myapp/src/server"
	"myapp/src/utils"
)

func main() {
    cfg,err := utils.LoadConfig()
    if err != nil {
        panic(err)
    }
    
    logger := logs.NewLogger(cfg.Log)

    dbConn, err := db.NewPostgresConnection(cfg.Database, logger)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer dbConn.Close()

    redisClient := redis.NewRedisClient(cfg.Redis, logger)
    defer redisClient.Close()

    app := &utils.App{
        Config: cfg,
        Logger: logger,
        DB:     dbConn,
        Redis:  redisClient,
    }

    srv := server.NewServer(app)
    srv.Start()
}

