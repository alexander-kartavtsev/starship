package converter

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func PartsFilterToRepo(filter model.PartsFilter) bson.M {
	repoFilter := bson.M{}

	var filters []bson.M

	if len(filter.Uuids) > 0 {
		var bUuids []bson.M
		for _, uuid := range filter.Uuids {
			bUuids = append(bUuids, bson.M{"uuid": uuid})
		}
		filters = append(filters, bson.M{"$or": bUuids})
	}

	if len(filter.Categories) > 0 {
		var bCat []bson.M
		for _, cat := range filter.Categories {
			bCat = append(bCat, bson.M{"categories": cat})
		}
		filters = append(filters, bson.M{"$or": bCat})
	}

	if len(filter.Names) > 0 {
		var bNames []bson.M
		for _, name := range filter.Names {
			bNames = append(bNames, bson.M{"name": bson.M{"$regex": name}})
		}
		filters = append(filters, bson.M{"$or": bNames})
	}

	if len(filter.ManufacturerCountries) > 0 {
		var bManCan []bson.M
		for _, can := range filter.ManufacturerCountries {
			bManCan = append(bManCan, bson.M{"manufacturer.country": bson.M{"$regex": can}})
		}
		filters = append(filters, bson.M{"$or": bManCan})
	}

	if len(filter.ManufacturerNames) > 0 {
		var bManNam []bson.M
		for _, name := range filter.ManufacturerNames {
			bManNam = append(bManNam, bson.M{"manufacturer.country": bson.M{"$regex": name}})
		}
		filters = append(filters, bson.M{"$or": bManNam})
	}

	if len(filter.Tags) > 0 {
		var bTags []bson.M
		for _, tag := range filter.Tags {
			bTags = append(bTags, bson.M{"tags": bson.M{"$regex": tag}})
		}
		filters = append(filters, bson.M{"$or": bTags})
	}

	if len(filters) > 0 {
		repoFilter = bson.M{
			"$and": filters,
		}
	}
	log.Printf("Собрали фильтр: %v\n", repoFilter)

	return repoFilter
}
