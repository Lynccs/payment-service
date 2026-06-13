package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Lynccs/payment-service/internal/app/dto"
	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/Lynccs/payment-service/internal/app/repo"
	appsvc "github.com/Lynccs/payment-service/internal/app/service"
)

type budgetService struct {
	budgetRepo repo.BudgetRepository
	walletRepo repo.WalletRepository
}

func NewBudgetService(budgetRepo repo.BudgetRepository, walletRepo repo.WalletRepository) appsvc.BudgetService {
	return &budgetService{budgetRepo: budgetRepo, walletRepo: walletRepo}
}

func (s *budgetService) ListBudgets(ctx context.Context, userID int, period time.Time) ([]dto.BudgetResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}

	items, err := s.budgetRepo.List(ctx, wallet.ID, period)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BudgetResponse, len(items))
	for i, b := range items {
		result[i] = budgetToResponse(b)
	}
	return result, nil
}

func (s *budgetService) CreateBudget(ctx context.Context, userID int, req dto.CreateBudgetRequest) (dto.BudgetResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.BudgetResponse{}, fmt.Errorf("get wallet: %w", err)
	}

	period, err := time.Parse("2006-01-02", req.Period)
	if err != nil {
		return dto.BudgetResponse{}, fmt.Errorf("invalid period: %w", err)
	}

	created, err := s.budgetRepo.Create(ctx, models.BudgetLimit{
		WalletID:   wallet.ID,
		CategoryID: req.CategoryID,
		Period:     period,
		Amount:     req.LimitAmount,
	})
	if err != nil {
		return dto.BudgetResponse{}, err
	}

	return dto.BudgetResponse{
		ID:          created.ID,
		CategoryID:  created.CategoryID,
		LimitAmount: created.Amount,
	}, nil
}

func (s *budgetService) UpdateBudget(ctx context.Context, userID int, budgetID int, req dto.UpdateBudgetRequest) (dto.BudgetResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.BudgetResponse{}, fmt.Errorf("get wallet: %w", err)
	}

	updated, err := s.budgetRepo.Update(ctx, models.BudgetLimit{
		ID:       budgetID,
		WalletID: wallet.ID,
		Amount:   req.LimitAmount,
	})
	if err != nil {
		return dto.BudgetResponse{}, err
	}

	return dto.BudgetResponse{
		ID:          updated.ID,
		CategoryID:  updated.CategoryID,
		LimitAmount: updated.Amount,
	}, nil
}

func (s *budgetService) DeleteBudget(ctx context.Context, userID int, budgetID int) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}
	return s.budgetRepo.Delete(ctx, budgetID, wallet.ID)
}

func budgetToResponse(b models.BudgetWithSpent) dto.BudgetResponse {
	return dto.BudgetResponse{
		ID:           b.ID,
		CategoryID:   b.CategoryID,
		CategoryName: b.CategoryName,
		CategoryIcon: b.CategoryIcon,
		LimitAmount:  b.Amount,
		Spent:        b.Spent,
	}
}
