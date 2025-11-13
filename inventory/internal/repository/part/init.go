package part

import (
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
)

// GetCollectionParts
// Функция для заполнения б/д тестовыми данными
func GetCollectionParts() []interface{} {
	parts := GetAllParts()
	var res []interface{}
	for _, part := range parts {
		res = append(res, part)
	}
	return res
}

func GetAllParts() map[string]repoModel.Part {
	parts := make(map[string]repoModel.Part)

	part := repoModel.Part{
		Uuid:     "a0ad507d-2b70-49e4-9378-3d92ebf9e405",
		Name:     "Двигатель",
		Category: repoModel.CATEGORY_ENGINE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Корпорация двигателей",
			Country: "Россия",
		},
		Tags:          []string{"двигатель", "Россия"},
		Price:         14250000,
		StockQuantity: 7,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "905ed12b-3934-45e1-a9af-67f00e00ff3d",
		Name:     "Ракетное топливо",
		Category: repoModel.CATEGORY_FUEL,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Дизель",
			Country: "Россия",
		},
		Tags:          []string{"топливо", "Россия", "ракета"},
		Price:         220,
		StockQuantity: 11365,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "21c03a7f-0760-4d10-86a4-3273c025a3c3",
		Name:     "Челночное топливо",
		Category: repoModel.CATEGORY_FUEL,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Chunga Changa",
			Country: "Китай",
		},
		Tags:          []string{"топливо", "Китай", "челнок"},
		Price:         330,
		StockQuantity: 24508,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "b10c92a3-d630-4aa6-9432-01167f77b20e",
		Name:     "Левое крыло",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "China wing",
			Country: "Китай",
		},
		Tags:          []string{"крыло", "Китай", "лев"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "7bfa48c8-50fa-42d8-8582-ba4d3ee410da",
		Name:     "Левое крыло",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Русское крыло",
			Country: "Россия",
		},
		Tags:          []string{"крыло", "Россия", "лев"},
		Price:         1848300,
		StockQuantity: 12,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "8aa70329-dfb9-4840-8566-73522b2a0dbf",
		Name:     "Правое крыло",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Русское крыло",
			Country: "Россия",
		},
		Tags:          []string{"крыло", "Россия", "прав"},
		Price:         1848300,
		StockQuantity: 11,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "4723d3ab-3650-4e30-9f0e-56cf4d1af44d",
		Name:     "Правое крыло",
		Category: repoModel.CATEGORY_WING,
		Manufacturer: repoModel.Manufacturer{
			Name:    "China wing",
			Country: "Китай",
		},
		Tags:          []string{"крыло", "Китай", "прав"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "1601a086-0973-4ea7-adac-357a96b6d8fa",
		Name:     "Иллюминатор круглый",
		Category: repoModel.CATEGORY_PORTHOLE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Окно в мир",
			Country: "Россия",
		},
		Tags:          []string{"окно", "Россия", "круг", "иллюминатор"},
		Price:         325000,
		StockQuantity: 84,
	}
	parts[part.Uuid] = part

	part = repoModel.Part{
		Uuid:     "8e04fd86-3cca-4500-9889-a910d3a5f1f9",
		Name:     "Иллюминатор квадратный",
		Category: repoModel.CATEGORY_PORTHOLE,
		Manufacturer: repoModel.Manufacturer{
			Name:    "Windows",
			Country: "Америка",
		},
		Tags:          []string{"окно", "Америка", "квадрат", "иллюминатор"},
		Price:         548000,
		StockQuantity: 14,
	}
	parts[part.Uuid] = part

	return parts
}
