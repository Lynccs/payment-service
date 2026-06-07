package postgres

import (
	"context"
	"fmt"

	"github.com/Lynccs/payment-service/internal/app/models"
	"github.com/Lynccs/payment-service/internal/app/repo"
	"github.com/jmoiron/sqlx"
)

type PaymentRepo struct {
	db *sqlx.DB
}

func NewPaymentRepo(db *sqlx.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(ctx context.Context, p models.Payment) (models.Payment, error) {
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO payments
			(wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date, created_at`,
		p.WalletID, p.IsIncome, p.Amount, p.CategoryID, p.PaymentMethodID, p.Description, p.TransactionDate,
	).StructScan(&p)
	if err != nil {
		return models.Payment{}, fmt.Errorf("create payment: %w", err)
	}
	return p, nil
}

func (r *PaymentRepo) List(ctx context.Context, filter repo.PaymentFilter, limit int) ([]models.PaymentListItem, error) {
	var payments []models.PaymentListItem
	err := r.db.SelectContext(ctx, &payments,
		`SELECT p.id, p.is_income, p.amount,
		        c.name  AS category_name,
		        m.name  AS method_name,
		        p.description, p.transaction_date
		FROM payments p
		LEFT JOIN payment_categories c ON c.id = p.category_id
		LEFT JOIN payment_methods     m ON m.id = p.payment_method_id
		WHERE p.wallet_id = $1 AND p.transaction_date BETWEEN $2 AND $3
		ORDER BY p.transaction_date DESC
		LIMIT $4`,
		filter.WalletID, filter.DateFrom, filter.DateTo, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepo) GetStats(ctx context.Context, filter repo.PaymentFilter) (repo.DashboardStats, error) {
	var stats repo.DashboardStats

	err := r.db.QueryRowxContext(ctx,
		`SELECT
			COALESCE(SUM(amount) FILTER (WHERE is_income = true), 0)  AS income,
			COALESCE(SUM(amount) FILTER (WHERE is_income = false), 0) AS expenses,
			COUNT(*) AS tx_count
		FROM payments
		WHERE wallet_id = $1 AND transaction_date BETWEEN $2 AND $3`,
		filter.WalletID, filter.DateFrom, filter.DateTo,
	).Scan(&stats.Income, &stats.Expenses, &stats.TxCount)
	if err != nil {
		return repo.DashboardStats{}, fmt.Errorf("get stats: %w", err)
	}

	rows, err := r.db.QueryxContext(ctx,
		`SELECT c.name AS category_name, SUM(p.amount) AS total
		FROM payments p
		JOIN payment_categories c ON c.id = p.category_id
		WHERE p.wallet_id = $1 AND p.is_income = false AND p.transaction_date BETWEEN $2 AND $3
		GROUP BY c.name
		ORDER BY total DESC`,
		filter.WalletID, filter.DateFrom, filter.DateTo,
	)
	if err != nil {
		return repo.DashboardStats{}, fmt.Errorf("get category stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cs repo.CategoryStat
		if err := rows.Scan(&cs.CategoryName, &cs.Total); err != nil {
			return repo.DashboardStats{}, err
		}
		stats.CategoryStats = append(stats.CategoryStats, cs)
	}

	return stats, nil
}

func (r *PaymentRepo) GetCategories(ctx context.Context, isIncome bool) ([]models.PaymentCategory, error) {
	var cats []models.PaymentCategory
	err := r.db.SelectContext(ctx, &cats,
		`SELECT p.id, p.name, p.icon, p.group_id,
		        g.name AS group_name, g.is_income
		FROM payment_categories p
		JOIN category_groups g ON g.id = p.group_id
		WHERE g.is_income = $1
		ORDER BY g.id, p.name`,
		isIncome,
	)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}
	return cats, nil
}

func (r *PaymentRepo) GetMethods(ctx context.Context) ([]models.PaymentMethod, error) {
	var methods []models.PaymentMethod
	err := r.db.SelectContext(ctx, &methods,
		`SELECT id, name FROM payment_methods ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("get methods: %w", err)
	}
	return methods, nil
}
