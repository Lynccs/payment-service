package dto

type BudgetResponse struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	CategoryIcon string  `json:"category_icon"`
	LimitAmount  float64 `json:"limit_amount"`
	Spent        float64 `json:"spent"`
}

type CreateBudgetRequest struct {
	CategoryID  int     `json:"category_id" binding:"required"`
	LimitAmount float64 `json:"limit_amount" binding:"required,gt=0"`
	Period      string  `json:"period" binding:"required"` // "YYYY-MM-01"
}

type UpdateBudgetRequest struct {
	LimitAmount float64 `json:"limit_amount" binding:"required,gt=0"`
}
