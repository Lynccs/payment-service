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
	//	TODO DataBase
	// TODO JWT
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8082"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env-default:"5s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"10s"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("cannot read .env file: %s", err)
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

	return &cfg
}
