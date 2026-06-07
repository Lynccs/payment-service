package models

import "time"

type Payment struct {
	ID              int       `db:"id"`
	WalletID        int       `db:"wallet_id"`
	IsIncome        bool      `db:"is_income"`
	Amount          float64   `db:"amount"`
	CategoryID      *int      `db:"category_id"`
	PaymentMethodID *int      `db:"payment_method_id"`
	Description     *string   `db:"description"`
	TransactionDate time.Time `db:"transaction_date"`
	CreatedAt       time.Time `db:"created_at"`
}

type PaymentCategory struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	IsIncome bool   `db:"is_income"`
}

type PaymentMethod struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
