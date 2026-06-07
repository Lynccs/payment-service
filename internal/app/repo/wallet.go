package repo

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/models"
)

type WalletRepository interface {
	Create(ctx context.Context, userID int) (models.Wallet, error)
	GetByUserID(ctx context.Context, userID int) (models.Wallet, error)
	UpdateBalance(ctx context.Context, walletID int, delta float64) error
}
