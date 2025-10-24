package model

type OrderInfo struct {
	UserUuid  string
	PartUuids []string
}

type OrderUpdateInfo struct {
	TransactionUuid string
	PartUuids       []string
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type Order struct {
	OrderUuid       string
	UserUuid        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUuid string
	PaymentMethod   PaymentMethod
	Status          OrderStatus
}

type PaymentMethod string

const (
	unknown       PaymentMethod = "UNKNOWN"
	Card          PaymentMethod = "CARD"
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
