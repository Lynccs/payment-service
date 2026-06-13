package service

import (
	"context"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
)

type BudgetService interface {
	ListBudgets(ctx context.Context, userID int, period time.Time) ([]dto.BudgetResponse, error)
	CreateBudget(ctx context.Context, userID int, req dto.CreateBudgetRequest) (dto.BudgetResponse, error)
	UpdateBudget(ctx context.Context, userID int, budgetID int, req dto.UpdateBudgetRequest) (dto.BudgetResponse, error)
	DeleteBudget(ctx context.Context, userID int, budgetID int) error
}
