package model

type PayOrderRequest struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}

func (r *PayOrderRequest) GetPaymentMethod() PaymentMethod {
	return r.PaymentMethod
}
