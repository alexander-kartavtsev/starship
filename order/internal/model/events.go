package model

type OrderKafkaEvent struct {
	Uuid            string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
	Type            string
}
