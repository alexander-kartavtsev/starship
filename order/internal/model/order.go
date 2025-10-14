package model

type OrderInfo struct {
}

type OrderUpdateInfo struct {
	Status OrderStatus
}

type Order struct {
	OrderUuid       string
	TotalPrice      float64
	PaymentMethod   PaymentMethod
	TransactionUuid string
	Status          OrderStatus
}

type OrderCreateRes struct {
	OrderUuid  string
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

type OrderStatus string

const (
	PendingPayment OrderStatus = "PENDING_PAYMENT"
	Paid           OrderStatus = "PAID"
	Cancelled      OrderStatus = "CANCELLED"
)
