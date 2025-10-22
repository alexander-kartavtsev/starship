package part

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/alexander-kartavtsev/starship/inventory/internal/repository"
	repoModel "github.com/alexander-kartavtsev/starship/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu sync.RWMutex
	// data       map[string]repoModel.Part
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Ошибка получения данных из mongoDb: %v\n", err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Ошибка при закрытии курсора: %v\n", cerr)
		}
	}()

	var val []repoModel.Part
	err = cursor.All(ctx, &val)
	if err != nil {
		log.Printf("Ошибка получения данных из mongoDb: %v\n", err)
	}
	if val == nil {
		_, err = collection.InsertMany(ctx, GetCollectionParts())
		if err != nil {
			log.Printf("Ошибка заполнения mongoDb: %v\n", err)
		}
	}

	return &repository{
		// data:       GetAllParts(),
		collection: collection,
	}
}
