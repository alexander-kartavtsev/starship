package converter

import (
	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func PartsToProto(parts map[string]model.Part) map[string]*inventoryV1.Part {
	protoParts := map[string]*inventoryV1.Part{}
	for partUuid, part := range parts {
		protoParts[partUuid] = PartToProto(&part)
	}
	return protoParts
}

func PartToProto(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return &inventoryV1.Part{}
	}
	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(part.Category),
		Dimensions:    DimentionsToProto(&part.Dimensions),
		Manufacturer:  ManufacturerToProto(&part.Manufacturer),
		Tags:          part.Tags,
	}
}

func DimentionsToProto(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return &inventoryV1.Dimensions{}
	}
	return &inventoryV1.Dimensions{
		Width:  dimensions.Width,
		Length: dimensions.Length,
		Weight: dimensions.Weight,
		Height: dimensions.Height,
	}
}

func ManufacturerToProto(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return &inventoryV1.Manufacturer{}
	}
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToModel(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		ManufacturerNames:     filter.ManufacturerNames,
		Tags:                  filter.Tags,
	}
}

func CategoriesToModel(categories []inventoryV1.Category) []model.Category {
	var res []model.Category
	for _, category := range categories {
		res = append(res, model.Category(category))
	}
	return res
}
