package handlers

import (
	"encoding/csv"
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

	limit := 10
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	page := 1
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	filter.Offset = (page - 1) * limit

	if cat := c.Query("category"); cat != "" {
		if v, err := strconv.Atoi(cat); err == nil {
			filter.CategoryID = &v
		}
	}

	if t := c.Query("type"); t == "income" || t == "expense" {
		isIncome := t == "income"
		filter.IsIncome = &isIncome
	}

	if m := c.Query("method"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			filter.MethodID = &v
		}
	}

	filter.Search = c.Query("search")
	filter.ExcludeSystem = c.Query("exclude_system") == "true"

	result, err := h.svc.GetTransactions(c.Request.Context(), userID, filter, limit)
	if err != nil {
		h.log.Error("get transactions", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) UpdateTransaction(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	txID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.svc.UpdateTransaction(c.Request.Context(), userID, txID, req)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		h.log.Error("update transaction", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) DeleteTransaction(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	txID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteTransaction(c.Request.Context(), userID, txID); err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		h.log.Error("delete transaction", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
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

func (h *PaymentHandler) ReconcileBalance(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)

	var req dto.ReconcileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.ReconcileBalance(c.Request.Context(), userID, req.TargetBalance); err != nil {
		h.log.Error("reconcile balance", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PaymentHandler) ExportTransactionsCSV(c *gin.Context) {
	userID := c.GetInt(middleware.UserIDKey)
	filter, ok := parseDateFilter(c)
	if !ok {
		return
	}

	if cat := c.Query("category"); cat != "" {
		if v, err := strconv.Atoi(cat); err == nil {
			filter.CategoryID = &v
		}
	}
	if t := c.Query("type"); t == "income" || t == "expense" {
		isIncome := t == "income"
		filter.IsIncome = &isIncome
	}
	if m := c.Query("method"); m != "" {
		if v, err := strconv.Atoi(m); err == nil {
			filter.MethodID = &v
		}
	}
	filter.Search = c.Query("search")

	result, err := h.svc.GetTransactions(c.Request.Context(), userID, filter, 100000)
	if err != nil {
		h.log.Error("export transactions csv", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="transactions.csv"`)

	// UTF-8 BOM for correct Excel rendering
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"Дата", "Категорія", "Тип транзакції", "Тип платіжу", "Опис", "Сума"})
	for _, tx := range result.Data {
		txType := "Витрата"
		if tx.IsIncome {
			txType = "Дохід"
		}
		cat, method, desc := "", "", ""
		if tx.CategoryName != nil {
			cat = *tx.CategoryName
		}
		if tx.MethodName != nil {
			method = *tx.MethodName
		}
		if tx.Description != nil {
			desc = *tx.Description
		}
		_ = w.Write([]string{
			tx.TransactionDate.Format("02.01.2006"),
			cat,
			txType,
			method,
			desc,
			strconv.FormatFloat(tx.Amount, 'f', 2, 64),
		})
	}
	w.Flush()
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
	default:
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	return repo.PaymentFilter{DateFrom: dateFrom, DateTo: now}, true
}
