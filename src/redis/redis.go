package redis

import (
	"myapp/src/utils"

	"github.com/go-redis/redis/v8"
)



func NewRedisClient(config utils.RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     config.Address,
        Password: config.Password,
        DB:       config.DB,
    })
}