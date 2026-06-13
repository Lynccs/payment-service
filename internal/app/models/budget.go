package models

import "time"

type BudgetLimit struct {
	ID         int       `db:"id"`
	WalletID   int       `db:"wallet_id"`
	CategoryID int       `db:"category_id"`
	Period     time.Time `db:"period"`
	Amount     float64   `db:"amount"`
}

type BudgetWithSpent struct {
	ID           int       `db:"id"`
	WalletID     int       `db:"wallet_id"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	CategoryIcon string    `db:"category_icon"`
	Period       time.Time `db:"period"`
	Amount       float64   `db:"amount"`
	Spent        float64   `db:"spent"`
}
