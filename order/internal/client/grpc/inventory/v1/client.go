package v1

import (
	def "github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
