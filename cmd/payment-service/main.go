package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lynccs/payment-service/internal/app"
	"github.com/Lynccs/payment-service/internal/app/service"
	"github.com/Lynccs/payment-service/internal/app/transport/http"
	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/Lynccs/payment-service/internal/pkg/logger"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
	"github.com/Lynccs/payment-service/internal/repository/postgres"
	"github.com/Lynccs/payment-service/internal/services"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)
	log.Info("initializing app", slog.String("env", cfg.Env))

	db, err := postgres.NewConnection(postgres.Config(cfg.Database))
	if err != nil {
		log.Error("Failed to connect to PostgreSQL", sl.Err(err))
		os.Exit(1)
	}
	defer postgres.Close(db)

	log.Info("Connected to PostgreSQL")

	userRepo := postgres.NewUserRepo(db)
	walletRepo := postgres.NewWalletRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)
	budgetRepo := postgres.NewBudgetRepo(db)

	svc := &service.Services{
		User:    services.NewUserService(userRepo, walletRepo),
		Payment: services.NewPaymentService(paymentRepo, walletRepo),
		Budget:  services.NewBudgetService(budgetRepo, walletRepo),
	}

	srv := http.NewServer(cfg, log, svc)

	a := app.New(cfg, log, srv)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := a.Run(ctx); err != nil {
		log.Error("app stopped with error", sl.Err(err))
		os.Exit(1)
	}

	log.Info("app stopped cleanly")
}
