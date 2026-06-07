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

type PaymentListItem struct {
	ID              int       `db:"id"`
	IsIncome        bool      `db:"is_income"`
	Amount          float64   `db:"amount"`
	CategoryName    *string   `db:"category_name"`
	MethodName      *string   `db:"method_name"`
	Description     *string   `db:"description"`
	TransactionDate time.Time `db:"transaction_date"`
}

type CategoryGroup struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	IsIncome bool   `db:"is_income"`
}

type PaymentCategory struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Icon      string `db:"icon"`
	GroupID   int    `db:"group_id"`
	GroupName string `db:"group_name"`
	IsIncome  bool   `db:"is_income"`
}

type PaymentMethod struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
