package converter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func PartsFilterToRepo(filter model.PartsFilter) bson.M {
	var filters []bson.M

	filters = prepareFilterIs("uuid", filter.Uuids, filters)
	filters = prepareFilterIs("categories", filter.Categories, filters)
	filters = prepareFilterIn("name", filter.Names, filters)
	filters = prepareFilterIn("manufacturer.country", filter.ManufacturerCountries, filters)
	filters = prepareFilterIn("manufacturer.name", filter.ManufacturerNames, filters)
	filters = prepareFilterIn("tags", filter.Tags, filters)

	return bindFilters(filters)
}

func bindFilters(filters []bson.M) bson.M {
	filter := bson.M{}
	if len(filters) > 0 {
		filter = bson.M{
			"$and": filters,
		}
	}
	return filter
}

func prepareFilterIs[T string | model.Category](field string, values []T, filters []bson.M) []bson.M {
	if len(values) > 0 {
		var resFilter []bson.M
		for _, value := range values {
			resFilter = append(resFilter, bson.M{field: value})
		}
		filters = append(filters, bson.M{"$or": resFilter})
	}
	return filters
}

func prepareFilterIn(field string, values []string, filters []bson.M) []bson.M {
	if len(values) > 0 {
		var resFilter []bson.M
		for _, value := range values {
			resFilter = append(resFilter, bson.M{field: bson.M{"$regex": value}})
		}
		filters = append(filters, bson.M{"$or": resFilter})
	}
	return filters
}
