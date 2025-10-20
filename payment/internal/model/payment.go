package model

type PayOrderRequest struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

type PaymentMethod int32

const (
	UNKNOWN PaymentMethod = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

func (r *PayOrderRequest) GetPaymentMethod() PaymentMethod {
	return r.PaymentMethod
}
