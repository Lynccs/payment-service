package handlers

import (
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

type PaymentHandler struct {
	svc service.PaymentService
	log *slog.Logger
}

func NewPaymentHandler(svc service.PaymentService, log *slog.Logger) *PaymentHandler {
	return &PaymentHandler{svc: svc, log: log}
}

func (h *PaymentHandler) CreateTransaction(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.CreateTransaction(c.Request.Context(), userID, req)
	if err != nil {
		h.log.Error("create transaction", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *PaymentHandler) GetTransactions(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)
	filter, ok := parseDateFilter(c)
	if !ok {
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	txs, err := h.svc.GetTransactions(c.Request.Context(), userID, filter, limit)
	if err != nil {
		h.log.Error("get transactions", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, txs)
}

func (h *PaymentHandler) GetDashboardStats(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)
	filter, ok := parseDateFilter(c)
	if !ok {
		return
	}

	stats, err := h.svc.GetDashboardStats(c.Request.Context(), userID, filter)
	if err != nil {
		h.log.Error("get dashboard stats", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *PaymentHandler) GetCategories(c *gin.Context) {
	isIncome := c.Query("is_income") == "true"

	cats, err := h.svc.GetCategories(c.Request.Context(), isIncome)
	if err != nil {
		h.log.Error("get categories", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, cats)
}

func (h *PaymentHandler) GetMethods(c *gin.Context) {
	methods, err := h.svc.GetMethods(c.Request.Context())
	if err != nil {
		h.log.Error("get methods", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, methods)
}

func parseDateFilter(c *gin.Context) (repo.PaymentFilter, bool) {
	now := time.Now()

	if from := c.Query("date_from"); from != "" {
		dateFrom, err := time.Parse("2006-01-02", from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from, expected YYYY-MM-DD"})
			return repo.PaymentFilter{}, false
		}
		dateTo := now
		if to := c.Query("date_to"); to != "" {
			dateTo, err = time.Parse("2006-01-02", to)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to, expected YYYY-MM-DD"})
				return repo.PaymentFilter{}, false
			}
			dateTo = dateTo.Add(24*time.Hour - time.Second)
		}
		return repo.PaymentFilter{DateFrom: dateFrom, DateTo: dateTo}, true
	}

	var dateFrom time.Time
	switch c.Query("preset") {
	case "week":
		dateFrom = now.AddDate(0, 0, -7)
	case "3months":
		dateFrom = now.AddDate(0, -3, 0)
	case "year":
		dateFrom = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default: // "month" і будь-що інше
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	return repo.PaymentFilter{DateFrom: dateFrom, DateTo: now}, true
}
