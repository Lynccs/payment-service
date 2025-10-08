package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lynccs/payment-service/internal/app"
	"github.com/Lynccs/payment-service/internal/app/transport/http"
	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/Lynccs/payment-service/internal/pkg/logger"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("initializing app", slog.String("env", cfg.Env))

	srv := http.NewServer(cfg, log)

	a := app.New(cfg, log, srv)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := a.Run(ctx); err != nil {
		log.Error("app stopped with error", sl.Err(err))
		os.Exit(1)
	}

	log.Info("app stopped cleanly")
}
