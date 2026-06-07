package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Lynccs/payment-service/internal/app/service"
	"github.com/Lynccs/payment-service/internal/app/transport/http/handlers"
	"github.com/Lynccs/payment-service/internal/pkg/config"
	"github.com/Lynccs/payment-service/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	log    *slog.Logger
	srv    *http.Server
}

func NewServer(cfg *config.Config, log *slog.Logger, services *service.Services) *Server {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(log))

	userHandler := handlers.NewUserHandler(services.User, log, cfg.JWT.Secret, cfg.JWT.TTL)

	r.Static("/static", "./web/static")
	r.StaticFile("/login", "./web/templates/login.html")
	r.StaticFile("/register", "./web/templates/register.html")

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(cfg.JWT.Secret))
		{
			protected.GET("/users/me", userHandler.GetMe)
		}
	}

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
