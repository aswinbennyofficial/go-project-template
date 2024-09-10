package db

import (
	"context"
	"fmt"
	"myapp/src/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)


type PostgresConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
}


func NewPostgresConnection(config utils.PostgresConfig) (*pgxpool.Pool, error) {
    connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)
    return pgxpool.Connect(context.Background(), connString)
}