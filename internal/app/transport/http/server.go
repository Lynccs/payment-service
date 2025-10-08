package http

import (
	"context"
	"fmt"
	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	log    *slog.Logger
	srv    *http.Server
}

func NewServer(cfg *config.Config, log *slog.Logger) *Server {
	r := gin.Default()

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &Server{
		router: r,
		cfg:    cfg,
		log:    log,
		srv:    srv,
	}
}

func (s *Server) Start() error {
	s.log.Info("starting http server", "address", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting http server: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down http server...")
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("error shutting down http server: %w", err)
	}
	return nil
}
