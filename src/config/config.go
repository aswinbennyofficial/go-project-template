package config

import (
	"fmt"
	"os"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type App struct {
    Config *Config
    Logger zerolog.Logger
    Postgres     *pgxpool.Pool
    Redis  *redis.Client
    Cassandra *gocql.Session
}


type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Postgres PostgresConfig `mapstructure:"postgres"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Log      LogConfig      `mapstructure:"log"`
    Auth     AuthConfig     `mapstructure:"auth"`
    Cassandra CassandraConfig `mapstructure:"cassandra"`
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

type CassandraConfig struct {
    Hosts    []string `mapstructure:"hosts"`
    Keyspace string   `mapstructure:"keyspace"`
    Port    int      `mapstructure:"port"`
    Username string `mapstructure:"username,omitempty"`
    Password string `mapstructure:"password,omitempty"`
    Consistency string `mapstructure:"consistency"`
    Replication CasandraReplication `mapstructure:"replication"`
    ProtoVersion int `mapstructure:"proto_version"`
    Migrations MigrationsConfig `mapstructure:"migrations"`
    IsKeySpaceSet bool `mapstructure:"keyspace_isset"`
}

type CasandraReplication struct{
    Strategy string `mapstructure:"strategy"`
    Factor int `mapstructure:"replication_factor"`
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

    // overwite log level from env if exists
    if os.Getenv("LOG_LEVEL") != "" {
        config.Log.Level = os.Getenv("LOG_LEVEL")
    }

    return &config, nil
}