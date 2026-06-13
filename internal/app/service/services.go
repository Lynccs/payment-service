package service

type Services struct {
	User    UserService
	Payment PaymentService
	Budget  BudgetService
}
