package redis

import (
    "context"
    "myapp/src/utils"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/rs/zerolog"
)

func NewRedisClient(config utils.RedisConfig, log zerolog.Logger) *redis.Client {
    var client *redis.Client
    var err error

    // Retry logic for Redis connection
    for i := 0; i < 30; i++ { // Retry up to 30 times (5 minutes total)
        client = redis.NewClient(&redis.Options{
            Addr:     config.Address,
            Password: config.Password,
            DB:       config.DB,
        })

        // Test the connection
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _, err = client.Ping(ctx).Result()
        cancel() 

        if err == nil {
            log.Info().Msg("Connected to Redis successfully")
            break
        }

        log.Warn().Err(err).Msgf("Failed to connect to Redis, retrying (%d/30)...", i+1)
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }

    if err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to Redis after retries")
        return nil
    }

    return client
}
