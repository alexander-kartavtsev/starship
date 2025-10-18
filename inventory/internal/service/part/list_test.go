package part

import (
	"errors"
	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
	"log"
)

func (s *ServiceSuite) TestService_List() {
	allTestParts := converter.PartsToModel(getAllParts())

	tests := []struct {
		name     string
		param    model.PartsFilter
		err      error
		expected map[string]model.Part
	}{
		{
			name: "inventory_service_List_By_Uuids_correct",
			param: model.PartsFilter{
				Uuids: []string{"part_uuid_1", "part_uuid_5", "part_uuid_20"},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_1": allTestParts["part_uuid_1"],
				"part_uuid_5": allTestParts["part_uuid_5"],
			},
		},
		{
			name: "inventory_service_List_By_Names_correct",
			param: model.PartsFilter{
				Names: []string{
					"Тестовое топливо 1",
					"Тестовое крыло левое 1",
					"Тестовое крыло правое 1",
				},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_3": allTestParts["part_uuid_3"],
				"part_uuid_5": allTestParts["part_uuid_5"],
				"part_uuid_8": allTestParts["part_uuid_8"],
			},
		},
		{
			name: "inventory_service_List_By_Categories_correct",
			param: model.PartsFilter{
				Categories: []model.Category{
					model.CATEGORY_ENGINE,
					model.CATEGORY_WING,
				},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_1": allTestParts["part_uuid_1"],
				"part_uuid_2": allTestParts["part_uuid_2"],
				"part_uuid_5": allTestParts["part_uuid_5"],
				"part_uuid_6": allTestParts["part_uuid_6"],
				"part_uuid_7": allTestParts["part_uuid_7"],
				"part_uuid_8": allTestParts["part_uuid_8"],
			},
		},
		{
			name: "inventory_service_List_By_Manufacturer_country_correct",
			param: model.PartsFilter{
				ManufacturerCountries: []string{
					"Вторая",
					"Третья",
				},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_4":  allTestParts["part_uuid_4"],
				"part_uuid_5":  allTestParts["part_uuid_5"],
				"part_uuid_8":  allTestParts["part_uuid_8"],
				"part_uuid_10": allTestParts["part_uuid_10"],
			},
		},
		{
			name: "inventory_service_List_By_Manufacturer_name_correct",
			param: model.PartsFilter{
				ManufacturerNames: []string{
					"Test corporation",
					"Окно",
					"Windows",
				},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_4":  allTestParts["part_uuid_4"],
				"part_uuid_5":  allTestParts["part_uuid_5"],
				"part_uuid_8":  allTestParts["part_uuid_8"],
				"part_uuid_9":  allTestParts["part_uuid_9"],
				"part_uuid_10": allTestParts["part_uuid_10"],
			},
		},
		{
			name: "inventory_service_List_By_Tags_correct",
			param: model.PartsFilter{
				Tags: []string{
					"прав",
					"круг",
				},
			},
			err: nil,
			expected: map[string]model.Part{
				"part_uuid_7": allTestParts["part_uuid_7"],
				"part_uuid_8": allTestParts["part_uuid_8"],
				"part_uuid_9": allTestParts["part_uuid_9"],
			},
		},
		{
			name:     "inventory_service_List_By_empty_filter_correct",
			param:    model.PartsFilter{},
			err:      nil,
			expected: allTestParts,
		},
		{
			name: "inventory_service_List_By_Uuids_Err_Not_found",
			param: model.PartsFilter{
				Uuids: []string{"part_uuid_not_exists"},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
		{
			name: "inventory_service_List_By_Names_Err_Not_found",
			param: model.PartsFilter{
				Names: []string{"name_not_exists"},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
		{
			name: "inventory_service_List_By_Categories_Err_Not_found",
			param: model.PartsFilter{
				Categories: []model.Category{model.CATEGORY_UNKNOWN},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
		{
			name: "inventory_service_List_By_Manuf_country_Err_Not_found",
			param: model.PartsFilter{
				ManufacturerCountries: []string{"country_not_exists"},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
		{
			name: "inventory_service_List_By_Manuf_name_Err_Not_found",
			param: model.PartsFilter{
				ManufacturerNames: []string{"name_not_exists"},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
		{
			name: "inventory_service_List_By_Tags_Err_Not_found",
			param: model.PartsFilter{
				Tags: []string{"tag_not_exists"},
			},
			err:      model.ErrPartNotFound,
			expected: map[string]model.Part{},
		},
	}
	for _, test := range tests {
		log.Println(test.name)
		s.inventoryRepository.On("List", s.ctx).Return(allTestParts, test.err).Once()
		res, err := s.service.List(s.ctx, test.param)
		s.Require().True(errors.Is(err, test.err))
		s.Require().Equal(test.expected, res)
	}
}

func (s *ServiceSuite) TestConverter_PartsToModel() {
	partUuid := "any_part_uuid"
	partsRepo := map[string]repoModel.Part{
		partUuid: repoModel.Part{
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
		partUuid: model.Part{
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

func getAllParts() map[string]repoModel.Part {
	parts := make(map[string]repoModel.Part)

	part := repoModel.Part{
		Uuid:     "part_uuid_1",
		Name:     "Тестовый двигатель 1",
		Category: repoModel.CATEGORY_ENGINE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Тестовая корпорация 1",
			Country: "Первая",
		},
		Tags:          []string{"двигатель", "Первая"},
		Price:         14250000,
		StockQuantity: 7,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_2",
		Name:     "Тестовый двигатель 2",
		Category: repoModel.CATEGORY_ENGINE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Тестовая корпорация 2",
			Country: "Первая",
		},
		Tags:          []string{"двигатель", "Первая"},
		Price:         14250000,
		StockQuantity: 7,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_3",
		Name:     "Тестовое топливо 1",
		Category: repoModel.CATEGORY_FUEL,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Дизель",
			Country: "Первая",
		},
		Tags:          []string{"топливо", "Первая", "type_one"},
		Price:         220,
		StockQuantity: 11365,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_4",
		Name:     "Тестовое топливо 2",
		Category: repoModel.CATEGORY_FUEL,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Test corporation 3",
			Country: "Вторая",
		},
		Tags:          []string{"топливо", "Вторая", "type_two"},
		Price:         330,
		StockQuantity: 24508,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_5",
		Name:     "Тестовое крыло левое 1",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Test corporation 2",
			Country: "Вторая",
		},
		Tags:          []string{"крыло", "Вторая", "лев"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_6",
		Name:     "Тестовое крыло левое 2",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Тестовое крыло",
			Country: "Первая",
		},
		Tags:          []string{"крыло", "Первая", "лев"},
		Price:         1848300,
		StockQuantity: 12,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_7",
		Name:     "Тестовое крыло правое 2",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Тестовое крыло",
			Country: "Первая",
		},
		Tags:          []string{"крыло", "Первая", "прав"},
		Price:         1848300,
		StockQuantity: 11,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_8",
		Name:     "Тестовое крыло правое 1",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Test corporation 2",
			Country: "Вторая",
		},
		Tags:          []string{"крыло", "Вторая", "прав"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_9",
		Name:     "Иллюминатор круглый",
		Category: repoModel.CATEGORY_PORTHOLE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Окно в мир",
			Country: "Первая",
		},
		Tags:          []string{"окно", "Первая", "круг", "иллюминатор"},
		Price:         325000,
		StockQuantity: 84,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "part_uuid_10",
		Name:     "Иллюминатор квадратный",
		Category: repoModel.CATEGORY_PORTHOLE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Windows",
			Country: "Третья",
		},
		Tags:          []string{"окно", "Третья", "квадрат", "иллюминатор"},
		Price:         548000,
		StockQuantity: 14,
	}
	parts[part.Uuid] = part

	return parts
}
