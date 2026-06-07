package dto

import "time"

type CreateTransactionRequest struct {
	IsIncome        bool    `json:"is_income"`
	Amount          float64 `json:"amount" binding:"required,gt=0"`
	CategoryID      *int    `json:"category_id"`
	PaymentMethodID *int    `json:"payment_method_id"`
	Description     *string `json:"description"`
	TransactionDate string  `json:"transaction_date" binding:"required"`
}

type TransactionResponse struct {
	ID              int        `json:"id"`
	IsIncome        bool       `json:"is_income"`
	Amount          float64    `json:"amount"`
	CategoryName    *string    `json:"category_name"`
	MethodName      *string    `json:"method_name"`
	Description     *string    `json:"description"`
	TransactionDate time.Time  `json:"transaction_date"`
}

type DashboardStatsResponse struct {
	Balance       float64         `json:"balance"`
	Income        float64         `json:"income"`
	Expenses      float64         `json:"expenses"`
	TxCount       int             `json:"tx_count"`
	CategoryStats []CategoryStat  `json:"category_stats"`
}

type CategoryStat struct {
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
}

type PaymentCategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	GroupName string `json:"group_name"`
	IsIncome  bool   `json:"is_income"`
}

type PaymentMethodResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
