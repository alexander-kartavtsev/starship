package payment

import def "github.com/alexander-kartavtsev/starship/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
