package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var part model.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		log.Printf("Ошибка получения part из коллекции: %v\n", err)
		return model.Part{}, model.ErrPartNotFound
	}
	return part, nil
}
