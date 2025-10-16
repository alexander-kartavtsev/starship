package converter

import (
	"github.com/alexander-kartavtsev/starship/payment/internal/model"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		UserUuid:      req.UserUuid,
		OrderUuid:     req.OrderUuid,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}
}
