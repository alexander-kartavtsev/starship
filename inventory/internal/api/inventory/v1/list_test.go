package v1

import (
	"errors"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func (a *ApiSuite) TestApi_ListPartsOk() {
	protoFilter := inventoryV1.PartsFilter{
		Uuids:                 []string{"any_part_uuid"},
		Names:                 []string{"any_names"},
		Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_FUEL},
		ManufacturerCountries: []string{"any_country"},
		ManufacturerNames:     []string{"any_name"},
		Tags:                  []string{"any_tag"},
	}
	modelFilter := model.PartsFilter{
		Uuids:                 []string{"any_part_uuid"},
		Names:                 []string{"any_names"},
		Categories:            []model.Category{model.CATEGORY_FUEL},
		ManufacturerCountries: []string{"any_country"},
		ManufacturerNames:     []string{"any_name"},
		Tags:                  []string{"any_tag"},
	}
	modelParts := map[string]model.Part{
		"any_uuid": {
			Uuid:          "any_uuid",
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
		},
	}
	protoParts := map[string]*inventoryV1.Part{
		"any_uuid": {
			Uuid:          "any_uuid",
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
		},
	}
	req := inventoryV1.ListPartsRequest{
		Filter: &protoFilter,
	}
	resp := inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}

	a.inventoryService.
		On("List", a.ctx, modelFilter).
		Return(modelParts, nil).
		Once()

	res, err := a.api.ListParts(a.ctx, &req)
	a.Assert().True(errors.Is(err, nil))
	a.Assert().Equal(&resp, res)
}

func (a *ApiSuite) TestApi_ListParts_NilFilter() {
	modelFilter := model.PartsFilter{
		Uuids:                 []string(nil),
		Names:                 []string(nil),
		Categories:            []model.Category(nil),
		ManufacturerCountries: []string(nil),
		ManufacturerNames:     []string(nil),
		Tags:                  []string(nil),
	}
	modelParts := map[string]model.Part{
		"any_uuid": {
			Uuid:          "any_uuid",
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
		},
	}
	protoParts := map[string]*inventoryV1.Part{
		"any_uuid": {
			Uuid:          "any_uuid",
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
		},
	}
	req := inventoryV1.ListPartsRequest{
		Filter: nil,
	}
	resp := inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}

	a.inventoryService.
		On("List", a.ctx, modelFilter).
		Return(modelParts, nil).
		Once()

	res, err := a.api.ListParts(a.ctx, &req)
	a.Assert().True(errors.Is(err, nil))
	a.Assert().Equal(&resp, res)
}

func (a *ApiSuite) TestApi_ListPartsErr() {
	protoFilter := inventoryV1.PartsFilter{
		Uuids:                 []string{"any_part_uuid"},
		Names:                 []string{"any_names"},
		Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_FUEL},
		ManufacturerCountries: []string{"any_country"},
		ManufacturerNames:     []string{"any_name"},
		Tags:                  []string{"any_tag"},
	}
	modelFilter := model.PartsFilter{
		Uuids:                 []string{"any_part_uuid"},
		Names:                 []string{"any_names"},
		Categories:            []model.Category{model.CATEGORY_FUEL},
		ManufacturerCountries: []string{"any_country"},
		ManufacturerNames:     []string{"any_name"},
		Tags:                  []string{"any_tag"},
	}
	req := inventoryV1.ListPartsRequest{
		Filter: &protoFilter,
	}

	a.inventoryService.
		On("List", a.ctx, modelFilter).
		Return(map[string]model.Part{}, model.ErrPartNotFound).
		Once()

	res, err := a.api.ListParts(a.ctx, &req)
	a.Assert().True(errors.Is(err, model.ErrPartNotFound))
	a.Assert().Equal(&inventoryV1.ListPartsResponse{}, res)
}
