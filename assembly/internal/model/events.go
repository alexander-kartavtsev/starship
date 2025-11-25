package model

type ShipAssembledKafkaEvent struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type OrderKafkaEvent struct {
	Uuid            string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
	Type            string
}

type PaymentMethod string

const (
	Unknown       PaymentMethod = "UNKNOWN"
	Card          PaymentMethod = "CARD"
	Sbp           PaymentMethod = "SBP"
	CreditCard    PaymentMethod = "CREDIT_CARD"
	InvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
