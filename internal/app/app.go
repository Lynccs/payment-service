package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
)

type App struct {
	config *config.Config
	logger *slog.Logger
	server Server
}

func New(cfg *config.Config, log *slog.Logger, server Server) *App {
	return &App{
		config: cfg,
		logger: log,
		server: server,
	}
}

func (a *App) Run(parentCtx context.Context) error {
	const op = "app.Run"
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.Start()
	}()

	var srvErr error
	select {
	case <-ctx.Done():
	case err := <-errCh:
		srvErr = err
		if err != nil {
			a.logger.Error("server exited with error", sl.Err(err))
		} else {
			a.logger.Info("server exited normally")
		}
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
	defer shutdownCancel()

	stopErr := a.server.Stop(shutdownCtx)
	if stopErr != nil {
		a.logger.Error("failed to stop server", sl.Err(stopErr))
	}

	if srvErr != nil {
		if stopErr != nil {
			return fmt.Errorf("%s: server error: %w; stop error: %v", op, srvErr, stopErr)
		}
		return fmt.Errorf("%s: server error: %w", op, srvErr)
	}
	if stopErr != nil {
		return fmt.Errorf("%s: failed to stop server: %w", op, stopErr)
	}

	a.logger.Info("graceful shutdown complete")
	return nil
}
