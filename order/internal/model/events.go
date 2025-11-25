package model

type OrderKafkaEvent struct {
	Uuid            string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
	Type            string
}

type ShipAssembledKafkaEvent struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}
