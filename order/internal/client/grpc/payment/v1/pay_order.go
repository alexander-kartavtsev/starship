package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/converter"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (c *client) PayOrder(ctx context.Context, req model.PayOrderRequest) (string, error) {
	response, err := c.generatedClient.PayOrder(ctx, converter.PayOrderRequestToProto(req))
	if err != nil {
		return "", err
	}
	// log.Printf("response in order grpc PayOrder: %v\n", response)
	return response.GetTransactionUuid(), nil
}
