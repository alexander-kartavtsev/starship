package part

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	def "github.com/alexander-kartavtsev/starship/inventory/internal/repository"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

const collectionName = "parts"

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu         sync.RWMutex
	collection *mongo.Collection
}

func NewRepository(_ context.Context, db *mongo.Database) *repository {
	return &repository{
		collection: db.Collection(collectionName),
	}
}

func (r *repository) InitFull(ctx context.Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := r.collection.FindOne(ctx, bson.M{})
	if res.Err() != nil {
		logger.Info(ctx, "Кажется, тут пусто")
		_, err := r.collection.InsertMany(ctx, GetCollectionParts())
		if err != nil {
			logger.Error(ctx, "Ошибка заполнения mongoDb", zap.Error(err))
			return err
		}
		logger.Info(ctx, "Заполнили тестовыми данными")
		return nil
	}
	logger.Info(ctx, "В коллекции есть данные", zap.Error(res.Err()))
	return nil
}
