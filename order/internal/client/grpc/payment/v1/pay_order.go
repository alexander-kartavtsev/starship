package v1

import (
	"context"
	"errors"

	"github.com/alexander-kartavtsev/starship/order/internal/converter"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (c *client) PayOrder(ctx context.Context, req model.PayOrderRequest) (string, error) {
	response, err := c.generatedClient.PayOrder(ctx, converter.PayOrderRequestToProto(req))
	if err != nil {
		return "", errors.New("ошибка соединения с сервисом payment")
	}
	return response.GetTransactionUuid(), nil
}
