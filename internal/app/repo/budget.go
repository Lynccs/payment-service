package repo

import (
	"context"
	"time"

	"github.com/Lynccs/payment-service/internal/app/models"
)

type BudgetRepository interface {
	List(ctx context.Context, walletID int, period time.Time) ([]models.BudgetWithSpent, error)
	Create(ctx context.Context, b models.BudgetLimit) (models.BudgetLimit, error)
	Update(ctx context.Context, b models.BudgetLimit) (models.BudgetLimit, error)
	Delete(ctx context.Context, id int, walletID int) error
}
