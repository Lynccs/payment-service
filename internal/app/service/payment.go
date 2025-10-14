package services

import (
	"context"
	"time"
)

type Payment struct {
	ID       string    `json:"id"`
	UserId   string    `json:"user_id"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

type PaymentService interface {
	GetAllPayment(ctx context.Context) ([]Payment, error)
	GetPaymentByID(ctx context.Context, id string) (*Payment, error)
	CreatePayment(ctx context.Context, payment *Payment) (*Payment, error)
}

func GetAllPayment(ctx context.Context) ([]Payment, error) {
	return nil, nil
}
