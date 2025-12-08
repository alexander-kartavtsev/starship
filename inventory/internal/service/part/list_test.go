package part

import (
	"errors"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
)

func (s *ServiceSuite) TestService_List_correct() {
	parts := map[string]model.Part{
		"any_uuid": {
			Uuid: "any_uuid",
		},
	}
	filter := model.PartsFilter{
		Uuids: []string{"any_uuid"},
	}
	ctx, span := tracing.StartSpan(s.ctx, "inventory.repository.List")
	span.End()

	s.inventoryRepository.
		On("List", ctx, filter).
		Return(parts, nil).
		Once()
	res, err := s.service.List(s.ctx, filter)
	s.Assert().Nil(err)
	s.Assert().Equal(res, parts)
}

func (s *ServiceSuite) TestService_List_err() {
	parts := map[string]model.Part{
		"any_uuid": {
			Uuid: "any_uuid",
		},
	}
	filter := model.PartsFilter{
		Uuids: []string{"any_uuid"},
	}
	ctx, span := tracing.StartSpan(s.ctx, "inventory.repository.List")
	span.End()

	s.inventoryRepository.
		On("List", ctx, filter).
		Return(parts, model.ErrPartListEmpty).
		Once()
	res, err := s.service.List(s.ctx, filter)
	s.Assert().True(errors.Is(err, model.ErrPartListEmpty))
	s.Assert().NotEqual(parts, res)
	s.Assert().Equal(map[string]model.Part{}, res)
}

func (s *ServiceSuite) TestConverter_PartsToModel() {
	partUuid := "any_part_uuid"
	partsRepo := map[string]repoModel.Part{
		partUuid: {
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
		},
	}

	expected := map[string]model.Part{
		partUuid: {
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
		},
	}

	res := converter.PartsToModel(partsRepo)
	s.Require().Equal(expected, res)
}
