package postgres

import (
	"context"
	"fmt"

	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type WalletRepo struct {
	db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) Create(ctx context.Context, userID int) (models.Wallet, error) {
	var w models.Wallet
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO wallets (user_id) VALUES ($1) RETURNING id, user_id, balance, created_at`,
		userID,
	).StructScan(&w)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("create wallet: %w", err)
	}
	return w, nil
}

func (r *WalletRepo) GetByUserID(ctx context.Context, userID int) (models.Wallet, error) {
	var w models.Wallet
	err := r.db.QueryRowxContext(ctx,
		`SELECT id, user_id, balance, created_at FROM wallets WHERE user_id = $1`,
		userID,
	).StructScan(&w)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("get wallet: %w", err)
	}
	return w, nil
}

func (r *WalletRepo) UpdateBalance(ctx context.Context, walletID int, delta float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`,
		delta, walletID,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	return nil
}
