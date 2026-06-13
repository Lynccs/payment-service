package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/repo"
	"github.com/Lynccs/payment-service/internal/app/service"
	"github.com/Lynccs/payment-service/internal/pkg/logger/sl"
	"github.com/Lynccs/payment-service/internal/pkg/middleware"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	svc service.BudgetService
	log *slog.Logger
}

func NewBudgetHandler(svc service.BudgetService, log *slog.Logger) *BudgetHandler {
	return &BudgetHandler{svc: svc, log: log}
}

func (h *BudgetHandler) ListBudgets(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	now := time.Now()
	period := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	if p := c.Query("period"); p != "" {
		if parsed, err := time.Parse("2006-01-02", p); err == nil {
			period = parsed
		}
	}

	budgets, err := h.svc.ListBudgets(c.Request.Context(), userID, period)
	if err != nil {
		h.log.Error("list budgets", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, budgets)
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	var req dto.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	budget, err := h.svc.CreateBudget(c.Request.Context(), userID, req)
	if err != nil {
		h.log.Error("create budget", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, budget)
}

func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	budget, err := h.svc.UpdateBudget(c.Request.Context(), userID, budgetID, req)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "budget not found"})
			return
		}
		h.log.Error("update budget", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	budgetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteBudget(c.Request.Context(), userID, budgetID); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "budget not found"})
			return
		}
		h.log.Error("delete budget", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
