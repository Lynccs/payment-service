package models

import "time"

type Payment struct {
	ID           string    `json:"id" db:"id"`
	FromWalletID string    `json:"from_wallet_id" db:"from_wallet_id"`
	ToWalletID   string    `json:"to_wallet_id" db:"to_wallet_id"`
	Amount       float64   `json:"amount" db:"amount"`
	Description  string    `json:"description" db:"description"`
	TypeID       string    `json:"type_id" db:"type_id"`
	StatusID     string    `json:"status_id" db:"status_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type PaymentType struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type PaymentStatus struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
