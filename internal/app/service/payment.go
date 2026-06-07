package service

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/repo"
)

type PaymentService interface {
	CreateTransaction(ctx context.Context, userID int, req dto.CreateTransactionRequest) (dto.TransactionResponse, error)
	GetTransactions(ctx context.Context, userID int, filter repo.PaymentFilter, limit int) ([]dto.TransactionResponse, error)
	GetDashboardStats(ctx context.Context, userID int, filter repo.PaymentFilter) (dto.DashboardStatsResponse, error)
	GetCategories(ctx context.Context, isIncome bool) ([]dto.PaymentCategoryResponse, error)
	GetMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error)
}
