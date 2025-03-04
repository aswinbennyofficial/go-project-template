// redis/client.go
package redis

import (
    "context"
    "myapp/src/config"
    "time"
    "github.com/redis/go-redis/v9"
    "github.com/rs/zerolog"
)

// NewRedisClientWrapper selects the appropriate Redis client based on config
func NewRedisClientWrapper(config config.RedisConfig, log zerolog.Logger) (redis.UniversalClient, error) {
    if config.Mode == "cluster" {
        return NewClusterRedisClient(config, log)
    }
    return NewStandaloneRedisClient(config, log)
}


// NewStandaloneRedisClient initializes a single Redis instance client
func NewStandaloneRedisClient(config config.RedisConfig, log zerolog.Logger) (*redis.Client, error) {
    var client *redis.Client
    var err error

    for i := 0; i < 30; i++ { // Retry up to 30 times (5 minutes total)
        client = redis.NewClient(&redis.Options{
            Addr:     config.Address,
            Username: config.Username,
            Password: config.Password,
            DB:       config.DB,
        })

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _, err = client.Ping(ctx).Result()
        cancel()

        if err == nil {
            log.Info().
                Str("address", config.Address).
                Int("db", config.DB).
                Bool("auth_enabled", config.Username != "" || config.Password != "").
                Msg("Connected to Redis (Standalone) successfully")
            return client, nil
        }

        log.Warn().
            Err(err).
            Int("attempt", i+1).
            Int("max_attempts", 30).
            Msg("Failed to connect to Redis (Standalone), retrying...")

        client.Close()
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }

    log.Fatal().
        Err(err).
        Str("address", config.Address).
        Msg("Failed to connect to Redis (Standalone) after all retries")
    return nil, err
}


// NewClusterRedisClient initializes a Redis Cluster client
func NewClusterRedisClient(config config.RedisConfig, log zerolog.Logger) (*redis.ClusterClient, error) {
    var client *redis.ClusterClient
    var err error

    for i := 0; i < 30; i++ { // Retry up to 30 times (5 minutes total)
        client = redis.NewClusterClient(&redis.ClusterOptions{
            Addrs:    config.ClusterAddresses,
            Username: config.Username,
            Password: config.Password,
        })

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        _, err = client.Ping(ctx).Result()
        cancel()

        if err == nil {
            log.Info().
                Strs("addresses", config.ClusterAddresses).
                Bool("auth_enabled", config.Username != "" || config.Password != "").
                Msg("Connected to Redis (Cluster) successfully")
            return client, nil
        }

        log.Warn().
            Err(err).
            Int("attempt", i+1).
            Int("max_attempts", 30).
            Msg("Failed to connect to Redis (Cluster), retrying...")

        client.Close()
        time.Sleep(10 * time.Second) // Wait 10 seconds before retrying
    }

    log.Fatal().
        Err(err).
        Strs("addresses", config.ClusterAddresses).
        Msg("Failed to connect to Redis (Cluster) after all retries")
    return nil, err
}


