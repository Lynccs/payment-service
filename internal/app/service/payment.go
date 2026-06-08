package service

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/repo"
)

type PaymentService interface {
	CreateTransaction(ctx context.Context, userID int, req dto.CreateTransactionRequest) (dto.TransactionResponse, error)
	UpdateTransaction(ctx context.Context, userID int, txID int, req dto.UpdateTransactionRequest) (dto.TransactionResponse, error)
	GetTransactions(ctx context.Context, userID int, filter repo.PaymentFilter, limit int) (dto.PaginatedTransactions, error)
	DeleteTransaction(ctx context.Context, userID int, txID int) error
	GetDashboardStats(ctx context.Context, userID int, filter repo.PaymentFilter) (dto.DashboardStatsResponse, error)
	ReconcileBalance(ctx context.Context, userID int, targetBalance float64) error
	GetCategories(ctx context.Context, isIncome bool) ([]dto.PaymentCategoryResponse, error)
	GetMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error)
}
