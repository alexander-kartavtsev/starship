package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alexander-kartavtsev/starship/inventory/internal/converter"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	log.Printf("%v\n", req)
	part, err := a.inventoryService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Запчасть с UUID %s не найдена", req.GetUuid())
	}

	log.Printf("%v\n", part)

	return &inventoryV1.GetPartResponse{
		Info: converter.PartToProto(part),
	}, nil
}
