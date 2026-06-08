package services

import (
	"context"
	"fmt"
	"math"
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

	return paymentToResponse(created), nil
}

func (s *paymentService) UpdateTransaction(ctx context.Context, userID int, txID int, req dto.UpdateTransactionRequest) (dto.TransactionResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("get wallet: %w", err)
	}

	txDate, err := time.Parse(time.RFC3339, req.TransactionDate)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("invalid transaction_date: %w", err)
	}

	p := models.Payment{
		ID:              txID,
		WalletID:        wallet.ID,
		IsIncome:        req.IsIncome,
		Amount:          req.Amount,
		CategoryID:      req.CategoryID,
		PaymentMethodID: req.PaymentMethodID,
		Description:     req.Description,
		TransactionDate: txDate,
	}

	updated, err := s.paymentRepo.Update(ctx, p)
	if err != nil {
		return dto.TransactionResponse{}, fmt.Errorf("update payment: %w", err)
	}

	return paymentToResponse(updated), nil
}

func (s *paymentService) GetTransactions(ctx context.Context, userID int, filter repo.PaymentFilter, limit int) (dto.PaginatedTransactions, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.PaginatedTransactions{}, fmt.Errorf("get wallet: %w", err)
	}
	filter.WalletID = wallet.ID

	total, err := s.paymentRepo.Count(ctx, filter)
	if err != nil {
		return dto.PaginatedTransactions{}, err
	}

	payments, err := s.paymentRepo.List(ctx, filter, limit)
	if err != nil {
		return dto.PaginatedTransactions{}, err
	}

	result := make([]dto.TransactionResponse, len(payments))
	for i, p := range payments {
		result[i] = listItemToResponse(p)
	}

	page := 1
	if limit > 0 && filter.Offset > 0 {
		page = filter.Offset/limit + 1
	}

	return dto.PaginatedTransactions{
		Data:  result,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *paymentService) DeleteTransaction(ctx context.Context, userID int, txID int) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	_, err = s.paymentRepo.Delete(ctx, txID, wallet.ID)
	if err != nil {
		return fmt.Errorf("delete payment: %w", err)
	}

	return nil
}

func (s *paymentService) GetDashboardStats(ctx context.Context, userID int, filter repo.PaymentFilter) (dto.DashboardStatsResponse, error) {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return dto.DashboardStatsResponse{}, fmt.Errorf("get wallet: %w", err)
	}
	filter.WalletID = wallet.ID

	balance, err := s.walletRepo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return dto.DashboardStatsResponse{}, fmt.Errorf("get balance: %w", err)
	}

	stats, err := s.paymentRepo.GetStats(ctx, filter)
	if err != nil {
		return dto.DashboardStatsResponse{}, err
	}

	catStats := make([]dto.CategoryStat, len(stats.CategoryStats))
	for i, c := range stats.CategoryStats {
		catStats[i] = dto.CategoryStat{CategoryName: c.CategoryName, Total: c.Total}
	}

	return dto.DashboardStatsResponse{
		Balance:       balance,
		Income:        stats.Income,
		Expenses:      stats.Expenses,
		TxCount:       stats.TxCount,
		CategoryStats: catStats,
	}, nil
}

func (s *paymentService) ReconcileBalance(ctx context.Context, userID int, targetBalance float64) error {
	wallet, err := s.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get wallet: %w", err)
	}

	currentBalance, err := s.walletRepo.GetBalance(ctx, wallet.ID)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}

	diff := targetBalance - currentBalance
	if math.Abs(diff) < 0.01 {
		return nil
	}

	isIncome := diff > 0
	amount := math.Abs(diff)

	cats, err := s.paymentRepo.GetCategories(ctx, isIncome)
	if err != nil {
		return fmt.Errorf("get categories: %w", err)
	}

	var catID *int
	for _, c := range cats {
		if c.Name == "Коригування балансу" {
			id := c.ID
			catID = &id
			break
		}
	}

	p := models.Payment{
		WalletID:        wallet.ID,
		IsIncome:        isIncome,
		Amount:          amount,
		CategoryID:      catID,
		TransactionDate: time.Now(),
	}

	_, err = s.paymentRepo.Create(ctx, p)
	return err
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
		CategoryID:      p.CategoryID,
		CategoryName:    p.CategoryName,
		CategoryIcon:    p.CategoryIcon,
		PaymentMethodID: p.PaymentMethodID,
		MethodName:      p.MethodName,
		Description:     p.Description,
		TransactionDate: p.TransactionDate,
	}
}
