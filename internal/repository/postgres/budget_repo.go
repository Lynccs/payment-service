package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/Lynccs/payment-service/internal/app/repo"
	"github.com/jmoiron/sqlx"
)

type BudgetRepo struct {
	db *sqlx.DB
}

func NewBudgetRepo(db *sqlx.DB) *BudgetRepo {
	return &BudgetRepo{db: db}
}

func (r *BudgetRepo) List(ctx context.Context, walletID int, period time.Time) ([]models.BudgetWithSpent, error) {
	var items []models.BudgetWithSpent
	err := r.db.SelectContext(ctx, &items, `
		SELECT
			bl.id,
			bl.wallet_id,
			bl.category_id,
			pc.name  AS category_name,
			pc.icon  AS category_icon,
			bl.period,
			bl.amount,
			COALESCE(SUM(p.amount), 0) AS spent
		FROM budget_limits bl
		JOIN payment_categories pc ON pc.id = bl.category_id
		LEFT JOIN payments p
			ON  p.category_id = bl.category_id
			AND p.wallet_id   = bl.wallet_id
			AND p.is_income   = false
			AND DATE_TRUNC('month', p.transaction_date) = bl.period
			AND p.deleted_at IS NULL
		WHERE bl.wallet_id = $1
		  AND bl.period = $2
		GROUP BY bl.id, bl.wallet_id, bl.category_id, pc.name, pc.icon, bl.period, bl.amount
		ORDER BY pc.name`,
		walletID, period,
	)
	if err != nil {
		return nil, fmt.Errorf("list budgets: %w", err)
	}
	return items, nil
}

func (r *BudgetRepo) Create(ctx context.Context, b models.BudgetLimit) (models.BudgetLimit, error) {
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO budget_limits (wallet_id, category_id, period, amount)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, wallet_id, category_id, period, amount`,
		b.WalletID, b.CategoryID, b.Period, b.Amount,
	).StructScan(&b)
	if err != nil {
		return models.BudgetLimit{}, fmt.Errorf("create budget: %w", err)
	}
	return b, nil
}

func (r *BudgetRepo) Update(ctx context.Context, b models.BudgetLimit) (models.BudgetLimit, error) {
	err := r.db.QueryRowxContext(ctx,
		`UPDATE budget_limits
		 SET amount = $1
		 WHERE id = $2 AND wallet_id = $3
		 RETURNING id, wallet_id, category_id, period, amount`,
		b.Amount, b.ID, b.WalletID,
	).StructScan(&b)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BudgetLimit{}, repo.ErrNotFound
		}
		return models.BudgetLimit{}, fmt.Errorf("update budget: %w", err)
	}
	return b, nil
}

func (r *BudgetRepo) Delete(ctx context.Context, id int, walletID int) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM budget_limits WHERE id = $1 AND wallet_id = $2`,
		id, walletID,
	)
	if err != nil {
		return fmt.Errorf("delete budget: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return repo.ErrNotFound
	}
	return nil
}
