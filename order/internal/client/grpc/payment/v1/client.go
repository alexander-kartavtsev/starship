package v1

import (
	def "github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
