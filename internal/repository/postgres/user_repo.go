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

var (
	ErrUserNotFound = fmt.Errorf("user not found")
)

// UserRepo реалізація UserRepository для PostgreSQL
type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) repo.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) (int, error) {
	op := "UserRepo.Create"
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (models.User, error) {
	op := "UserRepo.GetByID"
	query := `SELECT id, email, created_at FROM users WHERE id = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (models.User, error) {
	op := "UserRepo.GetByEmail"
	query := `SELECT id, email, created_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
