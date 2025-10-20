package part

import (
	"log"
	"strings"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func GetPartsByFilter(parts map[string]model.Part, filter *model.PartsFilter) (map[string]model.Part, error) {
	if filter == nil {
		return parts, model.ErrPartNotFound
	}

	filterUuids := filter.Uuids
	if filterUuids != nil {
		parts = getPartsByString(parts, filterUuids, "Uuid")
		log.Printf("parts: %v\n", parts)
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	filterNames := filter.Names
	if filterNames != nil {
		parts = getPartsByString(parts, filterNames, "Name")
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	filterCategories := filter.Categories
	if filterCategories != nil {
		parts = getPartsByCategories(parts, filterCategories)
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	filterCountries := filter.ManufacturerCountries
	if filterCountries != nil {
		parts = getPartsByCountry(parts, filterCountries, "Country")
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	filterManufacturer := filter.ManufacturerNames
	if filterManufacturer != nil {
		parts = getPartsByCountry(parts, filterManufacturer, "Name")
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	filterTags := filter.Tags
	if filterTags != nil {
		parts = getPartsByTags(parts, filterTags) // ???
		if len(parts) == 0 {
			return nil, model.ErrPartNotFound
		}
	}

	return parts, nil
}

// string start

func getPartsByString(parts map[string]model.Part, filter []string, name string) map[string]model.Part {
	res := map[string]model.Part{}

	for id, part := range parts {
		val, ok := GetStringFieldValue(&part, name)
		var check bool
		switch name {
		case "Uuid":
			check = checkStringValStrong(val, filter)
		default:
			check = checkStringVal(val, filter)
		}
		if ok && check {
			res[id] = part
		}
	}
	return res
}

func GetStringFieldValue(part *model.Part, name string) (string, bool) {
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

func checkStringValStrong(value string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	for _, filterValue := range filter {
		if strings.EqualFold(strings.ToLower(value), strings.ToLower(filterValue)) {
			return true
		}
	}
	return false
}

// string end
// ---
// categories start

func getPartsByCategories(parts map[string]model.Part, filter []model.Category) map[string]model.Part {
	res := map[string]model.Part{}

	// log.Println("  Это срез категорий")
	for id, part := range parts {
		val, ok := GetCategoryFieldValue(&part)
		// log.Printf("Получили значение поля Category: %s", val)
		if ok && checkCategoryVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetCategoryFieldValue(part *model.Part) (model.Category, bool) {
	return part.Category, true
}

func checkCategoryVal(value model.Category, filter []model.Category) bool {
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

func getPartsByCountry(parts map[string]model.Part, filter []string, field string) map[string]model.Part {
	res := map[string]model.Part{}

	for id, part := range parts {
		val, ok := GetCountryFieldValue(&part, field)
		if ok && checkStringVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetCountryFieldValue(part *model.Part, field string) (string, bool) {
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

func getPartsByTags(parts map[string]model.Part, filter []string) map[string]model.Part {
	res := map[string]model.Part{}

	// log.Println("  Это срез строк")
	for id, part := range parts {
		val, ok := GetTagsFieldValue(&part)
		// log.Printf("Получили значение поля %s: %s", name, val)
		if ok && checkTagsVal(val, filter) {
			res[id] = part
		}
	}
	return res
}

func GetTagsFieldValue(part *model.Part) ([]string, bool) {
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
