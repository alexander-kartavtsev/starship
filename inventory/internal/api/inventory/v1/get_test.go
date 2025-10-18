package v1

import (
	"errors"
	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *ApiSuite) TestApi_GetPartOk() {
	requestUuid := "ani_part_uuid"

	modelPart := model.Part{
		Uuid:          requestUuid,
		Name:          "any_name",
		Description:   "any_description",
		Price:         123.45,
		StockQuantity: 15,
		Category:      model.CATEGORY_PORTHOLE,
		Dimensions: model.Dimensions{
			Length: 123.45,
			Width:  456.12,
			Height: 123.45,
			Weight: 456.12,
		},
		Manufacturer: model.Manufacturer{
			Name:    "any manufacturer name",
			Country: "any country",
			Website: "any.web.site",
		},
		Tags: []string{"any_tag"},
	}
	protoPart := inventoryV1.Part{
		Uuid:          requestUuid,
		Name:          "any_name",
		Description:   "any_description",
		Price:         123.45,
		StockQuantity: 15,
		Category:      inventoryV1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 123.45,
			Width:  456.12,
			Height: 123.45,
			Weight: 456.12,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "any manufacturer name",
			Country: "any country",
			Website: "any.web.site",
		},
		Tags: []string{"any_tag"},
	}

	a.inventoryService.
		On("Get", a.ctx, requestUuid).
		Return(modelPart, nil).
		Once()
	res, err := a.api.GetPart(a.ctx, &inventoryV1.GetPartRequest{Uuid: requestUuid})
	a.Assert().True(errors.Is(err, nil))
	a.Assert().Equal(res.GetInfo(), &protoPart)
}

func (a *ApiSuite) TestApi_GetPartErr() {
	requestUuid := "any_uuid_part_not_found"
	modelPart := model.Part{}

	a.inventoryService.
		On("Get", a.ctx, requestUuid).
		Return(modelPart, model.ErrPartNotFound).
		Once()

	res, err := a.api.GetPart(a.ctx, &inventoryV1.GetPartRequest{Uuid: requestUuid})
	a.Assert().True(errors.Is(err, status.Errorf(codes.NotFound, "Запчасть с UUID %s не найдена", requestUuid)))
	a.Assert().Equal(&inventoryV1.GetPartResponse{}, res)
}
