package utils

import (
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/go-redis/redis/v8"
    "github.com/rs/zerolog"
)

type App struct {
    Config *Config
    Logger zerolog.Logger
    DB     *pgxpool.Pool
    Redis  *redis.Client
}