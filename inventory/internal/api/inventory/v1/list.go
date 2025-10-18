package v1

import (
	"context"

	"github.com/alexander-kartavtsev/starship/inventory/internal/converter"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	reqFilter := req.GetFilter()

	parts, err := a.inventoryService.List(ctx, converter.PartsFilterToModel(reqFilter))
	if err != nil {
		return &inventoryV1.ListPartsResponse{}, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
