package service

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/models"
)


type PaymentService interface {
	GetAllPayments(ctx context.Context) ([]models.Payment, error)
	GetPaymentByID(ctx context.Context, id string) (*models.Payment, error)
	CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
}

