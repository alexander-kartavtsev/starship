package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUuid := gofakeit.UUID()
	now := time.Now()
	part := bson.M{
		"uuid":        partUuid,
		"category":    1,
		"createdAt":   primitive.NewDateTimeFromTime(now),
		"description": gofakeit.ProductDescription(),
		"dimensions": bson.M{
			"length": 123.45,
			"width":  123.45,
			"height": 123.45,
			"weight": 123.45,
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
		},
		"name":          gofakeit.Product(),
		"price":         gofakeit.Price(10000, 10000000),
		"stockQuantity": gofakeit.Int(),
		"tags":          []string{"tag1", "tag2"},
	}

	// Используем базу данных из переменной окружения MONGO_INITDB_DATABASE
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return partUuid, nil
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
