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
		`INSERT INTO wallets (user_id) VALUES ($1) RETURNING id, user_id, initial_balance, created_at`,
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
		`SELECT id, user_id, initial_balance, created_at FROM wallets WHERE user_id = $1`,
		userID,
	).StructScan(&w)
	if err != nil {
		return models.Wallet{}, fmt.Errorf("get wallet: %w", err)
	}
	return w, nil
}

func (r *WalletRepo) GetBalance(ctx context.Context, walletID int) (float64, error) {
	var balance float64
	err := r.db.QueryRowxContext(ctx,
		`SELECT w.initial_balance + COALESCE(
			SUM(CASE WHEN p.is_income THEN p.amount ELSE -p.amount END), 0
		)
		FROM wallets w
		LEFT JOIN payments p ON p.wallet_id = w.id AND p.deleted_at IS NULL
		WHERE w.id = $1
		GROUP BY w.initial_balance`,
		walletID,
	).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}
	return balance, nil
}

func (r *WalletRepo) SetInitialBalance(ctx context.Context, walletID int, amount float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wallets SET initial_balance = $1 WHERE id = $2`,
		amount, walletID,
	)
	if err != nil {
		return fmt.Errorf("set initial balance: %w", err)
	}
	return nil
}
