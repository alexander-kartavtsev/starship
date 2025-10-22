package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func TestPartToProto(t *testing.T) {
	t.Run("converter_part_to_proto_err", func(t *testing.T) {
		var part *model.Part
		res := PartToProto(part)
		assert.Equal(t, &inventoryV1.Part{}, res)
	})
}

func TestDimensionsToProto(t *testing.T) {
	t.Run("converter_dimensions_to_proto_err", func(t *testing.T) {
		var dimensions *model.Dimensions
		res := DimensionsToProto(dimensions)
		assert.Equal(t, &inventoryV1.Dimensions{}, res)
	})
}

func TestManufacturerToProto(t *testing.T) {
	t.Run("converter_manufacturer_to_proto_err", func(t *testing.T) {
		var manufacturer *model.Manufacturer
		res := ManufacturerToProto(manufacturer)
		assert.Equal(t, &inventoryV1.Manufacturer{}, res)
	})
}
