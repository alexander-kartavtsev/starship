package part

import (
	"context"
	"errors"
	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
	repoMocks "github.com/alexander-kartavtsev/starship/inventory/internal/repository/mocks"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_Get_old(t *testing.T) {
	type inventoryServiceMockFunk func(t *testing.T) *repoMocks.InventoryRepository

	tests := []struct {
		name                     string
		param                    string
		err                      error
		expected                 model.Part
		inventoryServiceMockFunk inventoryServiceMockFunk
	}{
		{
			name:  "inventory_service_Get_correct",
			param: "part_uuid_1",
			err:   nil,
			expected: model.Part{
				Uuid: "part_uuid_1",
			},
			inventoryServiceMockFunk: func(t *testing.T) *repoMocks.InventoryRepository {
				partUuid := "part_uuid_1"
				mockRepo := repoMocks.NewInventoryRepository(t)
				mockRepo.
					On("Get", context.Background(), partUuid).
					Return(model.Part{Uuid: partUuid}, nil).
					Once()
				return mockRepo
			},
		},
		{
			name:     "inventory_service_Get_not_found",
			param:    "part_uuid_1",
			err:      model.ErrPartNotFound,
			expected: model.Part{},
			inventoryServiceMockFunk: func(t *testing.T) *repoMocks.InventoryRepository {
				partUuid := "part_uuid_1"
				mockRepo := repoMocks.NewInventoryRepository(t)
				mockRepo.
					On("Get", context.Background(), partUuid).
					Return(model.Part{}, model.ErrPartNotFound).
					Once()
				return mockRepo
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serv := NewService(test.inventoryServiceMockFunk(t))
			res, err := serv.Get(context.Background(), test.param)
			require.True(t, errors.Is(err, test.err))
			require.Equal(t, test.expected, res)
		})
	}
}

func TestConverter_PartToModel_old(t *testing.T) {
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
	require.Equal(t, expected, res)
}
