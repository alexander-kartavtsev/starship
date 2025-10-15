package v1

import (
	"context"

	clientConverter "github.com/alexander-kartavtsev/starship/order/internal/client/converter"
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) (map[string]model.Part, error) {
	// log.Printf("filter in order grpc ListParts: %v\n", filter)
	response, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	// log.Printf("response in order grpc ListParts: %v\n", response)
	return clientConverter.PartListToModel(response.GetParts()), nil
}
