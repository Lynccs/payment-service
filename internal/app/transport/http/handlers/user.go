package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/service"
	"github.com/Lynccs/payment-service/internal/pkg/jwt"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
	"github.com/Lynccs/payment-service/internal/pkg/middleware"
	"github.com/Lynccs/payment-service/internal/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc       service.UserService
	log       *slog.Logger
	jwtSecret string
	jwtTTL    time.Duration
}

func NewUserHandler(svc service.UserService, log *slog.Logger, jwtSecret string, jwtTTL time.Duration) *UserHandler {
	return &UserHandler{
		svc:       svc,
		log:       log,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request body", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req)
	if err != nil {
		h.log.Error("failed to register user", sl.Err(err))
		switch err {
		case services.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		case services.ErrPasswordTooShort:
			c.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	token, err := jwt.GenerateToken(user.ID, h.jwtSecret, h.jwtTTL)
	if err != nil {
		h.log.Error("failed to generate token", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Info("user registered", slog.Int("user_id", user.ID))
	c.JSON(http.StatusCreated, dto.AuthResponse{Token: token, User: user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request body", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		h.log.Error("login failed", sl.Err(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, h.jwtSecret, h.jwtTTL)
	if err != nil {
		h.log.Error("failed to generate token", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Info("user logged in", slog.Int("user_id", user.ID))
	c.JSON(http.StatusOK, dto.AuthResponse{Token: token, User: user})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.svc.GetProfile(c.Request.Context(), userID.(int))
	if err != nil {
		h.log.Error("failed to get profile", sl.Err(err))
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
