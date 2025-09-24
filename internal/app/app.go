package app

import (
	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/Lynccs/payment-service/internal/pkg/logger"
	"log/slog"
)

type App struct {
	config *config.Config
	logger *slog.Logger
}

func New() (*App, error) {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("initializing app", slog.String("env", cfg.Env))

	app := &App{
		config: cfg,
		logger: log,
	}

	return app, nil
}
