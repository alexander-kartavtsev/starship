package converter

import (
	"github.com/alexander-kartavtsev/starship/order/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		ManufacturerNames:     filter.ManufacturerNames,
		Tags:                  filter.Tags,
	}
}

func PartListToModel(parts map[string]*inventoryV1.Part) map[string]model.Part {
	res := map[string]model.Part{}
	for uuid, part := range parts {
		res[uuid] = PartToModel(part)
	}

	return res
}

func PartToModel(part *inventoryV1.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
	}
}
