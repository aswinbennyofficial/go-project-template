// redis/client.go
package redis

import (
    "context"
    "myapp/src/config"
    "time"
    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog"
)

func NewRedisClient(config config.RedisConfig, log zerolog.Logger) *redis.Client {
    var client *redis.Client
    var err error

    // Retry logic for Redis connection
    for i := 0; i < 30; i++ { // Retry up to 30 times (5 minutes total)
        if config.Username == "" {
            // No username, use only password if provided
            client = redis.NewClient(&redis.Options{
                Addr:     config.Address,
                Password: config.Password,
                DB:       config.DB,
            })
        } else if config.Password == "" {
            // Username provided, but no password
            client = redis.NewClient(&redis.Options{
                Addr:     config.Address,
                DB:       config.DB,
            })
        } else {
            // Username and password provided
            client = redis.NewClient(&redis.Options{
                Addr:     config.Address,
                Username: config.Username,
                Password: config.Password,
                DB:       config.DB,
            })
        }

        // Test the connection
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _, err = client.Ping(ctx).Result()
        cancel()

        if err == nil {
            log.Info().
                Str("address", config.Address).
                Int("db", config.DB).
                Bool("auth_enabled", config.Username != "" || config.Password != "").
                Msg("Connected to Redis successfully")
            break
        }

        log.Warn().
            Err(err).
            Int("attempt", i+1).
            Int("max_attempts", 30).
            Msg("Failed to connect to Redis, retrying...")

        if client != nil {
            _ = client.Close()
        }
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }

    if err != nil {
        log.Fatal().
            Err(err).
            Str("address", config.Address).
            Msg("Failed to connect to Redis after all retries")
        return nil
    }

    return client
}