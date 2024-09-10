package main

import (
    "myapp/src/db"
    "myapp/src/log"
    "myapp/src/redis"
    "myapp/src/server"
    "myapp/src/utils"
)

func main() {
    cfg,err := utils.LoadConfig()
    if err != nil {
        panic(err)
    }
    
    logger := log.NewLogger(cfg.Log)

    dbConn, err := db.NewPostgresConnection(cfg.Database)
    if err != nil {
        logger.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer dbConn.Close()

    redisClient := redis.NewRedisClient(cfg.Redis)
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