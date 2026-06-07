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

type paymentService struct {
	paymentRepo repo.PaymentRepository
	walletRepo  repo.WalletRepository
}

func NewPaymentService(paymentRepo repo.PaymentRepository, walletRepo repo.WalletRepository) appsvc.PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		walletRepo:  walletRepo,
	}
}

func (s *paymentService) CreateTransaction(ctx context.Context, userID int, req dto.CreateTransactionRequest) (dto.TransactionResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("get wallet: %w", err)
	}

	txDate, err := time.Parse(time.RFC3339, req.TransactionDate)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("invalid transaction_date: %w", err)
	}

	p := models.Payment{
		WalletID:        wallet.ID,
		IsIncome:        req.IsIncome,
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		PaymentMethodID: req.PaymentMethodID,
		Description:     req.Description,
		TransactionDate: txDate,
	}

	created, err := s.paymentRepo.Create(ctx, p)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("create payment: %w", err)
	}

	delta := created.Amount
	if !created.IsIncome {
		delta = -delta
	}
	if err := s.walletRepo.UpdateBalance(ctx, wallet.ID, delta); err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("update balance: %w", err)
	}

	return paymentToResponse(created), nil
}

func (s *paymentService) GetTransactions(ctx context.Context, userID int, filter repo.PaymentFilter, limit int) ([]dto.TransactionResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	filter.WalletID = wallet.ID

	payments, err := s.paymentRepo.List(ctx, filter, limit)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionResponse, len(payments))
	for i, p := range payments {
		result[i] = listItemToResponse(p)
	}
	return result, nil
}

func (s *paymentService) GetDashboardStats(ctx context.Context, userID int, filter repo.PaymentFilter) (dto.DashboardStatsResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.DashboardStatsResponse{}, fmt.Errorf("get wallet: %w", err)
	}
	filter.WalletID = wallet.ID

	stats, err := s.paymentRepo.GetStats(ctx, filter)
	if err != nil {
		return dto.DashboardStatsResponse{}, err
	}

	catStats := make([]dto.CategoryStat, len(stats.CategoryStats))
	for i, c := range stats.CategoryStats {
		catStats[i] = dto.CategoryStat{CategoryName: c.CategoryName, Total: c.Total}
	}

	return dto.DashboardStatsResponse{
		Balance:       wallet.Balance,
		Income:        stats.Income,
		Expenses:      stats.Expenses,
		TxCount:       stats.TxCount,
		CategoryStats: catStats,
	}, nil
}

func (s *paymentService) GetCategories(ctx context.Context, isIncome bool) ([]dto.PaymentCategoryResponse, error) {
	cats, err := s.paymentRepo.GetCategories(ctx, isIncome)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PaymentCategoryResponse, len(cats))
	for i, c := range cats {
		result[i] = dto.PaymentCategoryResponse{
			ID:        c.ID,
			Name:      c.Name,
			Icon:      c.Icon,
			GroupName: c.GroupName,
			IsIncome:  c.IsIncome,
		}
	}
	return result, nil
}

func (s *paymentService) GetMethods(ctx context.Context) ([]dto.PaymentMethodResponse, error) {
	methods, err := s.paymentRepo.GetMethods(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PaymentMethodResponse, len(methods))
	for i, m := range methods {
		result[i] = dto.PaymentMethodResponse{ID: m.ID, Name: m.Name}
	}
	return result, nil
}

func paymentToResponse(p models.Payment) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:              p.ID,
		IsIncome:        p.IsIncome,
		Amount:          p.Amount,
		Description:     p.Description,
		TransactionDate: p.TransactionDate,
	}
}

func listItemToResponse(p models.PaymentListItem) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:              p.ID,
		IsIncome:        p.IsIncome,
		Amount:          p.Amount,
		CategoryName:    p.CategoryName,
		MethodName:      p.MethodName,
		Description:     p.Description,
		TransactionDate: p.TransactionDate,
	}
}

