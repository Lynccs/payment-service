package repo

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/models"
)

type WalletRepository interface {
	Create(ctx context.Context, userID int) (models.Wallet, error)
	GetByUserID(ctx context.Context, userID int) (models.Wallet, error)
	GetBalance(ctx context.Context, walletID int) (float64, error)
	SetInitialBalance(ctx context.Context, walletID int, amount float64) error
}
