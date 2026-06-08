package repo

import (
	"context"
	"errors"
	"time"

	"github.com/Lynccs/payment-service/internal/app/models"
)

var ErrNotFound = errors.New("not found")

type PaymentFilter struct {
	WalletID      int
	DateFrom      time.Time
	DateTo        time.Time
	CategoryID    *int
	IsIncome      *bool
	MethodID      *int
	Search        string
	Offset        int
	ExcludeSystem bool
}

type DashboardStats struct {
	Income        float64
	Expenses      float64
	TxCount       int
	CategoryStats []CategoryStat
}

type CategoryStat struct {
	CategoryName string
	Total        float64
}

type PaymentRepository interface {
	Create(ctx context.Context, p models.Payment) (models.Payment, error)
	Update(ctx context.Context, p models.Payment) (models.Payment, error)
	List(ctx context.Context, filter PaymentFilter, limit int) ([]models.PaymentListItem, error)
	Count(ctx context.Context, filter PaymentFilter) (int, error)
	Delete(ctx context.Context, paymentID int, walletID int) (models.Payment, error)
	GetStats(ctx context.Context, filter PaymentFilter) (DashboardStats, error)
	GetCategories(ctx context.Context, isIncome bool) ([]models.PaymentCategory, error)
	GetMethods(ctx context.Context) ([]models.PaymentMethod, error)
}
