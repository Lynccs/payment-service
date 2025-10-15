package services

import (
	"context"

	"github.com/Lynccs/payment-service/internal/app/models"
	appsvc "github.com/Lynccs/payment-service/internal/app/service"
)

type paymentService struct {
	// repo repository.PaymentRepository
}

func NewPaymentService( /* repo repository.PaymentRepository */ ) appsvc.PaymentService {
	return &paymentService{ /* repo: repo */ }
}

// CreatePayment implements PaymentService.
func (s *paymentService) CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	panic("unimplemented")
}

// GetPaymentByID implements PaymentService.
func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (*models.Payment, error) {
	panic("unimplemented")
}


func (s *paymentService) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	// Implement logic to get all payments
	return []models.Payment{}, nil
}