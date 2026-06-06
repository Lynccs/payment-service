package dto

import (
	"time"

	"github.com/Lynccs/payment-service/internal/app/models"
)

type PaymentResponse struct {
	ID           string    `json:"id"`
	FromWalletID string    `json:"from_wallet_id"`
	ToWalletID   string    `json:"to_wallet_id"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToPaymentResponse(payment *models.Payment, paymentType string, paymentStatus string) *PaymentResponse {
	return &PaymentResponse{
		ID:           payment.ID,
		FromWalletID: payment.FromWalletID,
		ToWalletID:   payment.ToWalletID,
		Amount:       payment.Amount,
		Description:  payment.Description,
		Type:         paymentType,
		Status:       paymentStatus,
		CreatedAt:    payment.CreatedAt,
	}
}
