package repo

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/models"
)

// UserRepository інтерфейс для роботи з користувачами в БД
type UserRepository interface {
	Create(ctx context.Context, user models.User) (int, error)

	GetByID(ctx context.Context, id int) (models.User, error)

	GetByEmail(ctx context.Context, email string) (models.User, error)
}
