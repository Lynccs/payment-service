package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Database   Database `yaml:"database"`
	// TODO JWT
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"0.0.0.0:8082"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

type Database struct {
	Host            string        `env:"DB_HOST" env-required:"true"`
	Port            int           `env:"DB_PORT" env-default:"5432"`
	User            string        `env:"DB_USER" env-required:"true"`
	Password        string        `env:"DB_PASSWORD" env-required:"true"`
	Name            string        `env:"DB_NAME" env-required:"true"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"100"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"25"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"30m"`
}

func MustLoad() *Config {
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("cannot read .env file: %s", err)
		}
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set.")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not exists: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	if cfg.Database.Host == "" {
		log.Fatal("DB_HOST is required")
	}
	if cfg.Database.User == "" {
		log.Fatal("DB_USER is required")
	}
	if cfg.Database.Password == "" {
		log.Fatal("DB_PASSWORD is required")
	}
	if cfg.Database.Name == "" {
		log.Fatal("DB_NAME is required")
	}

	return &cfg
}
