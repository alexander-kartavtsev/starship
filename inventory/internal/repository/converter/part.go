package converter

import (
	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
)

func PartsToModel(parts map[string]repoModel.Part) map[string]model.Part {
	modelParts := map[string]model.Part{}
	for partUuid, part := range parts {
		modelParts[partUuid] = PartToModel(part)
	}
	return modelParts
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    DimensionsToModel(&part.Dimensions),
		Manufacturer:  ManufacturerToModel(&part.Manufacturer),
		Tags:          part.Tags,
	}
}

func DimensionsToModel(dimensions *repoModel.Dimensions) model.Dimensions {
	if dimensions == nil {
		return model.Dimensions{}
	}
	return model.Dimensions{
		Width:  dimensions.Width,
		Length: dimensions.Length,
		Weight: dimensions.Weight,
		Height: dimensions.Height,
	}
}

func ManufacturerToModel(manufacturer *repoModel.Manufacturer) model.Manufacturer {
	if manufacturer == nil {
		return model.Manufacturer{}
	}
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}
