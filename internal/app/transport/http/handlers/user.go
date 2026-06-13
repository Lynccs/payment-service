package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/service"
	appjwt "github.com/Lynccs/payment-service/internal/pkg/jwt"
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

func (h *UserHandler) setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", token, int(h.jwtTTL.Seconds()), "/", "", false, true)
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
		case services.ErrPasswordTooShort,
			services.ErrPasswordNoUpper,
			services.ErrPasswordNoDigit,
			services.ErrPasswordNoSymbol:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	token, err := appjwt.GenerateToken(user.ID, h.jwtSecret, h.jwtTTL)
	if err != nil {
		h.log.Error("failed to generate token", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.setAuthCookie(c, token)
	h.log.Info("user registered", slog.Int("user_id", user.ID))
	c.JSON(http.StatusCreated, user)
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

	token, err := appjwt.GenerateToken(user.ID, h.jwtSecret, h.jwtTTL)
	if err != nil {
		h.log.Error("failed to generate token", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.setAuthCookie(c, token)
	h.log.Info("user logged in", slog.Int("user_id", user.ID))
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
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
