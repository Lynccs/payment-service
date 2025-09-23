package app

import "github.com/Lynccs/payment-service/internal/pkg/config"

type App struct {
	config *config.Config
}

func New() (*App, error) {
	cfg := config.MustLoad()

	app := &App{config: cfg}

	return app, nil
}
