package utils

import (
    "gopkg.in/yaml.v3"
    "io/ioutil"
)

type Config struct {
    App      AppConfig      `yaml:"app"`
    Database PostgresConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Log      LogConfig      `yaml:"log"`
    Auth     AuthConfig     `yaml:"auth"`
}

type AppConfig struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`
    Port    int    `yaml:"port"`
}

type PostgresConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
}

type RedisConfig struct {
    Address  string `yaml:"address"`
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
}

type LogConfig struct {
    Level  string   `yaml:"level"`
    Output string   `yaml:"output"`
    File   FileLog  `yaml:"file"`
}

type FileLog struct {
    Path       string `yaml:"path"`
    MaxSize    int    `yaml:"max_size"`
    MaxAge     int    `yaml:"max_age"`
    MaxBackups int    `yaml:"max_backups"`
}

type AuthConfig struct {
    JWTSecret string `yaml:"jwt_secret"`
}

func LoadConfig() *Config {
    data, err := ioutil.ReadFile("config/app.yaml")
    if err != nil {
        panic(err)
    }

    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }

    return &config
}
