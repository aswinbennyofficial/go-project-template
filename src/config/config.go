package config

import (
    "fmt"
    "github.com/spf13/viper"
    "github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type App struct {
    Config *Config
    Logger zerolog.Logger
    DB     *pgxpool.Pool
    Redis  *redis.Client
}


type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Postgres PostgresConfig `mapstructure:"postgres"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Log      LogConfig      `mapstructure:"log"`
    Auth     AuthConfig     `mapstructure:"auth"`
}

type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
}

type PostgresConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    Migrations MigrationsConfig `mapstructure:"migrations"`
}

type MigrationsConfig struct {
    Enabled bool   `mapstructure:"enabled"`
    Path    string `mapstructure:"path"`
}

type RedisConfig struct {
    Address  string `mapstructure:"address"`
    Username string `mapstructure:"username,omitempty"`
    Password string `mapstructure:"password,omitempty"`
    DB       int    `mapstructure:"db"`
}

type LogConfig struct {
    Level  string  `mapstructure:"level"`
    Output string  `mapstructure:"output"`
    File   FileLog `mapstructure:"file"`
}

type FileLog struct {
    Path       string `mapstructure:"path"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxAge     int    `mapstructure:"max_age"`
    MaxBackups int    `mapstructure:"max_backups"`
}

type AuthConfig struct {
    JWTSecret string `mapstructure:"jwt_secret"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }

    return &config, nil
}