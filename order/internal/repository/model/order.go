package model

type OrderInfo struct {
}

type OrderUpdateInfo struct {
}

type Order struct {
	Uuid       string
	TotalPrice float64
}

type PaymentMethod string

const (
	unknown       PaymentMethod = "UNKNOWN"
	card          PaymentMethod = "CARD"
	sbp           PaymentMethod = "SBP"
	creditCard    PaymentMethod = "CREDIT_CARD"
	investorMoney PaymentMethod = "INVESTOR_MONEY"
)
