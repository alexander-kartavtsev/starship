package main

import (
	"log"
	"strings"

	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

func GetPartsByFilter(parts map[string]*inventoryV1.Part, filter *inventoryV1.PartsFilter) map[string]*inventoryV1.Part {
	if filter == nil {
		log.Println("Значение = nil")
		return parts
	}
	log.Printf("%v\n", filter)

	filterUuids := filter.Uuids
	if filterUuids != nil {
		parts = getPartsByString(parts, filterUuids, "Uuid")
		// log.Printf("%v\n", parts)
		if parts == nil {
			return nil
		}
	}
	// log.Printf("1: %v\n", parts)

	filterNames := filter.Names
	if filterNames != nil {
		parts = getPartsByString(parts, filterNames, "Name")
		if parts == nil {
			return nil
		}
	}
	// log.Printf("2: %v\n", parts)

	filterCategories := filter.Categories
	if filterCategories != nil {
		parts = getPartsByCategories(parts, filterCategories)
		if parts == nil {
			return nil
		}
	}
	// log.Printf("3: %v\n", parts)

	filterCountries := filter.ManufacturerCountries
	if filterCountries != nil {
		parts = getPartsByCountry(parts, filterCountries, "Country")
		if parts == nil {
			return nil
		}
	}
	// log.Printf("4: %v\n", parts)

	filterManufacturer := filter.ManufacturerNames
	if filterManufacturer != nil {
		parts = getPartsByCountry(parts, filterManufacturer, "Name")
		if parts == nil {
			return nil
		}
	}
	// log.Printf("5: %v\n", parts)

	filterTags := filter.Tags
	if filterTags != nil {
		parts = getPartsByTags(parts, filterTags) // ???
		if parts == nil {
			return nil
		}
	}
	// log.Printf("6: %v\n", parts)

	return parts
}

// string start

func getPartsByString(parts map[string]*inventoryV1.Part, filter []string, name string) map[string]*inventoryV1.Part {
	res := map[string]*inventoryV1.Part{}

	// log.Println("  Это срез строк")
	for id, part := range parts {
		val, ok := GetStringFieldValue(part, name)
		// log.Printf("Получили значение поля %s: %s", name, val)
		if ok && checkStringVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetStringFieldValue(part *inventoryV1.Part, name string) (string, bool) {
	switch name {
	case "Uuid":
		return part.Uuid, true
	case "Name":
		return part.Name, true
	default:
		return "", false
	}
}

func checkStringVal(value string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	for _, filterValue := range filter {
		if strings.Contains(strings.ToLower(value), strings.ToLower(filterValue)) {
			return true
		}
	}
	return false
}

// string end
// ---
// categories start

func getPartsByCategories(parts map[string]*inventoryV1.Part, filter []inventoryV1.Category) map[string]*inventoryV1.Part {
	res := map[string]*inventoryV1.Part{}

	// log.Println("  Это срез категорий")
	for id, part := range parts {
		val, ok := GetCategoryFieldValue(part)
		// log.Printf("Получили значение поля Category: %s", val)
		if ok && checkCategoryVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetCategoryFieldValue(part *inventoryV1.Part) (inventoryV1.Category, bool) {
	return part.Category, true
}

func checkCategoryVal(value inventoryV1.Category, filter []inventoryV1.Category) bool {
	if len(filter) == 0 {
		return true
	}
	for _, filterValue := range filter {
		if value == filterValue {
			return true
		}
	}
	return false
}

// categories end
// ---
// country start

func getPartsByCountry(parts map[string]*inventoryV1.Part, filter []string, field string) map[string]*inventoryV1.Part {
	res := map[string]*inventoryV1.Part{}

	// log.Println("  Это срез строк")
	for id, part := range parts {
		val, ok := GetCountryFieldValue(part, field)
		log.Printf("Получили значение поля %s: %s", field, val)
		if ok && checkStringVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetCountryFieldValue(part *inventoryV1.Part, field string) (string, bool) {
	switch field {
	case "Country":
		return part.Manufacturer.Country, true
	case "Name":
		return part.Manufacturer.Name, true
	default:
		return "", false
	}
}

// country end
// ---
// tags start

func getPartsByTags(parts map[string]*inventoryV1.Part, filter []string) map[string]*inventoryV1.Part {
	res := map[string]*inventoryV1.Part{}

	// log.Println("  Это срез строк")
	for id, part := range parts {
		val, ok := GetTagsFieldValue(part)
		// log.Printf("Получили значение поля %s: %s", name, val)
		if ok && checkTagsVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetTagsFieldValue(part *inventoryV1.Part) ([]string, bool) {
	return part.Tags, true
}

func checkTagsVal(tags, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	for _, filterValue := range filter {
		for _, value := range tags {
			if strings.Contains(strings.ToLower(value), strings.ToLower(filterValue)) {
				return true
			}
		}
	}
	return false
}

// tags end

func GetAllParts() map[string]*inventoryV1.Part {
	parts := make(map[string]*inventoryV1.Part)

	part := &inventoryV1.Part{
		Uuid:     "a0ad507d-2b70-49e4-9378-3d92ebf9e405",
		Name:     "Двигатель",
		Category: inventoryV1.Category_CATEGORY_ENGINE,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Корпорация двигателей",
			Country: "Россия",
		},
		Tags:          []string{"двигатель", "Россия"},
		Price:         14250000,
		StockQuantity: 7,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "905ed12b-3934-45e1-a9af-67f00e00ff3d",
		Name:     "Ракетное топливо",
		Category: inventoryV1.Category_CATEGORY_FUEL,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Дизель",
			Country: "Россия",
		},
		Tags:          []string{"топливо", "Россия", "ракета"},
		Price:         220,
		StockQuantity: 11365,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "21c03a7f-0760-4d10-86a4-3273c025a3c3",
		Name:     "Челночное топливо",
		Category: inventoryV1.Category_CATEGORY_FUEL,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Chunga Changa",
			Country: "Китай",
		},
		Tags:          []string{"топливо", "Китай", "челнок"},
		Price:         330,
		StockQuantity: 24508,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "b10c92a3-d630-4aa6-9432-01167f77b20e",
		Name:     "Левое крыло",
		Category: inventoryV1.Category_CATEGORY_WING,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "China wing",
			Country: "Китай",
		},
		Tags:          []string{"крыло", "Китай", "лев"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "7bfa48c8-50fa-42d8-8582-ba4d3ee410da",
		Name:     "Левое крыло",
		Category: inventoryV1.Category_CATEGORY_WING,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Русское крыло",
			Country: "Россия",
		},
		Tags:          []string{"крыло", "Россия", "лев"},
		Price:         1848300,
		StockQuantity: 12,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "8aa70329-dfb9-4840-8566-73522b2a0dbf",
		Name:     "Правое крыло",
		Category: inventoryV1.Category_CATEGORY_WING,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Русское крыло",
			Country: "Россия",
		},
		Tags:          []string{"крыло", "Россия", "прав"},
		Price:         1848300,
		StockQuantity: 11,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "4723d3ab-3650-4e30-9f0e-56cf4d1af44d",
		Name:     "Правое крыло",
		Category: inventoryV1.Category_CATEGORY_WING,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "China wing",
			Country: "Китай",
		},
		Tags:          []string{"крыло", "Китай", "прав"},
		Price:         2360800,
		StockQuantity: 4,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "1601a086-0973-4ea7-adac-357a96b6d8fa",
		Name:     "Иллюминатор круглый",
		Category: inventoryV1.Category_CATEGORY_PORTHOLE,
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Окно в мир",
			Country: "Россия",
		},
		Tags:          []string{"окно", "Россия", "круг", "иллюминатор"},
		Price:         325000,
		StockQuantity: 84,
	}
	parts[part.Uuid] = part

	part = &inventoryV1.Part{
		Uuid:     "8e04fd86-3cca-4500-9889-a910d3a5f1f9",
		Name:     "Иллюминатор квадратный",
		Category: inventoryV1.Category_CATEGORY_PORTHOLE,
		Manufacturer: &inventoryV1.Manufacturer{
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
