package part

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
)

func (s *ServiceSuite) TestService_Get() {
	tests := []struct {
		name     string
		param    string
		err      error
		expected model.Part
	}{
		{
			name:  "inventory_service_Get_correct",
			param: "part_uuid_1",
			err:   nil,
			expected: model.Part{
				Uuid: "part_uuid_1",
			},
		},
		{
			name:     "inventory_service_Get_not_found",
			param:    "part_uuid_1",
			err:      model.ErrPartNotFound,
			expected: model.Part{},
		},
	}

	for _, test := range tests {
		log.Println(test.name)
		s.inventoryRepository.On("Get", s.ctx, test.param).Return(test.expected, test.err).Once()
		res, err := s.service.Get(s.ctx, test.param)
		s.Assert().True(errors.Is(err, test.err))
		s.Assert().Equal(test.expected, res)
	}
}

func (s *ServiceSuite) TestConverter_PartToModel() {
	partUuid := "any_part_uuid"
	partRepo := repoModel.Part{
		Uuid:        partUuid,
		Name:        "any_name",
		Description: "any description",
		Price:       120,
		Category:    repoModel.CATEGORY_WING,
		Dimensions: repoModel.Dimensions{
			Length: 123,
			Width:  456,
			Height: 123,
			Weight: 456,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "any manufacturer",
			Country: "any country",
			Website: "anw.site.test",
		},
		Tags: []string{"any_tag"},
	}

	expected := model.Part{
		Uuid:        partUuid,
		Name:        "any_name",
		Description: "any description",
		Price:       120,
		Category:    model.CATEGORY_WING,
		Dimensions: model.Dimensions{
			Length: 123,
			Width:  456,
			Height: 123,
			Weight: 456,
		},
		Manufacturer: model.Manufacturer{
			Name:    "any manufacturer",
			Country: "any country",
			Website: "anw.site.test",
		},
		Tags: []string{"any_tag"},
	}

	res := converter.PartToModel(partRepo)
	s.Assert().Equal(expected, res)
}

func TestConverter_DimensionsToModel(t *testing.T) {
	modelDimensions := model.Dimensions{
		Width:  123.45,
		Length: 123.45,
		Weight: 123.45,
		Height: 123.45,
	}

	repoModelDimensions := repoModel.Dimensions{
		Width:  123.45,
		Length: 123.45,
		Weight: 123.45,
		Height: 123.45,
	}

	t.Run("dimensions_to_model_ok", func(t *testing.T) {
		res := converter.DimensionsToModel(&repoModelDimensions)
		assert.Equal(t, modelDimensions, res)
	})
	t.Run("dimensions_to_model_empty", func(t *testing.T) {
		res := converter.DimensionsToModel(&repoModel.Dimensions{})
		assert.Equal(t, model.Dimensions{}, res)
	})
}
