package postgres

import (
	"context"
	"database/sql"
	"errors"
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

func (r *PaymentRepo) Update(ctx context.Context, p models.Payment) (models.Payment, error) {
	err := r.db.QueryRowxContext(ctx,
		`UPDATE payments
		 SET is_income = $1, amount = $2, category_id = $3, payment_method_id = $4,
		     description = $5, transaction_date = $6
		 WHERE id = $7 AND wallet_id = $8 AND deleted_at IS NULL
		 RETURNING id, wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date, created_at`,
		p.IsIncome, p.Amount, p.CategoryID, p.PaymentMethodID, p.Description, p.TransactionDate,
		p.ID, p.WalletID,
	).StructScan(&p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payment{}, repo.ErrNotFound
		}
		return models.Payment{}, fmt.Errorf("update payment: %w", err)
	}
	return p, nil
}

func (r *PaymentRepo) List(ctx context.Context, filter repo.PaymentFilter, limit int) ([]models.PaymentListItem, error) {
	query := `SELECT p.id, p.is_income, p.amount,
	                 p.category_id, c.name AS category_name, c.icon AS category_icon,
	                 p.payment_method_id, m.name AS method_name,
	                 p.description, p.transaction_date
	          FROM payments p
	          LEFT JOIN payment_categories c ON c.id = p.category_id
	          LEFT JOIN category_groups g    ON g.id = c.group_id
	          LEFT JOIN payment_methods m    ON m.id = p.payment_method_id
	          WHERE p.wallet_id = $1
	            AND p.transaction_date BETWEEN $2 AND $3
	            AND p.deleted_at IS NULL`

	if filter.ExcludeSystem {
		query += ` AND (g.is_system IS NULL OR g.is_system = false)`
	}

	args := []interface{}{filter.WalletID, filter.DateFrom, filter.DateTo}
	n := 4

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", n)
		args = append(args, *filter.CategoryID)
		n++
	}
	if filter.IsIncome != nil {
		query += fmt.Sprintf(" AND p.is_income = $%d", n)
		args = append(args, *filter.IsIncome)
		n++
	}
	if filter.MethodID != nil {
		query += fmt.Sprintf(" AND p.payment_method_id = $%d", n)
		args = append(args, *filter.MethodID)
		n++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND p.description ILIKE $%d", n)
		args = append(args, "%"+filter.Search+"%")
		n++
	}

	query += fmt.Sprintf(" ORDER BY p.transaction_date DESC LIMIT $%d OFFSET $%d", n, n+1)
	args = append(args, limit, filter.Offset)

	var payments []models.PaymentListItem
	if err := r.db.SelectContext(ctx, &payments, query, args...); err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepo) Count(ctx context.Context, filter repo.PaymentFilter) (int, error) {
	query := `SELECT COUNT(*)
	          FROM payments p
	          WHERE p.wallet_id = $1
	            AND p.transaction_date BETWEEN $2 AND $3
	            AND p.deleted_at IS NULL`

	args := []interface{}{filter.WalletID, filter.DateFrom, filter.DateTo}
	n := 4

	if filter.CategoryID != nil {
		query += fmt.Sprintf(" AND p.category_id = $%d", n)
		args = append(args, *filter.CategoryID)
		n++
	}
	if filter.IsIncome != nil {
		query += fmt.Sprintf(" AND p.is_income = $%d", n)
		args = append(args, *filter.IsIncome)
		n++
	}
	if filter.MethodID != nil {
		query += fmt.Sprintf(" AND p.payment_method_id = $%d", n)
		args = append(args, *filter.MethodID)
		n++
	}
	if filter.Search != "" {
		query += fmt.Sprintf(" AND p.description ILIKE $%d", n)
		args = append(args, "%"+filter.Search+"%")
		n++
	}
	_ = n

	var count int
	if err := r.db.QueryRowxContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count payments: %w", err)
	}
	return count, nil
}

func (r *PaymentRepo) Delete(ctx context.Context, paymentID int, walletID int) (models.Payment, error) {
	var p models.Payment
	err := r.db.QueryRowxContext(ctx,
		`UPDATE payments
		 SET deleted_at = NOW()
		 WHERE id = $1 AND wallet_id = $2 AND deleted_at IS NULL
		 RETURNING id, wallet_id, is_income, amount, category_id, payment_method_id, description, transaction_date, created_at`,
		paymentID, walletID,
	).StructScan(&p)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payment{}, repo.ErrNotFound
		}
		return models.Payment{}, fmt.Errorf("delete payment: %w", err)
	}
	return p, nil
}

func (r *PaymentRepo) GetStats(ctx context.Context, filter repo.PaymentFilter) (repo.DashboardStats, error) {
	var stats repo.DashboardStats

	err := r.db.QueryRowxContext(ctx,
		`SELECT
			COALESCE(SUM(amount) FILTER (WHERE is_income = true), 0)  AS income,
			COALESCE(SUM(amount) FILTER (WHERE is_income = false), 0) AS expenses,
			COUNT(*) AS tx_count
		FROM payments
		WHERE wallet_id = $1
		  AND transaction_date BETWEEN $2 AND $3
		  AND deleted_at IS NULL`,
		filter.WalletID, filter.DateFrom, filter.DateTo,
	).Scan(&stats.Income, &stats.Expenses, &stats.TxCount)
	if err != nil {
		return repo.DashboardStats{}, fmt.Errorf("get stats: %w", err)
	}

	rows, err := r.db.QueryxContext(ctx,
		`SELECT c.name AS category_name, SUM(p.amount) AS total
		FROM payments p
		JOIN payment_categories c ON c.id = p.category_id
		JOIN category_groups g ON g.id = c.group_id
		WHERE p.wallet_id = $1
		  AND p.is_income = false
		  AND p.transaction_date BETWEEN $2 AND $3
		  AND p.deleted_at IS NULL
		  AND g.is_system = false
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
