package part

import (
	"context"
	"log"

	"github.com/alexander-kartavtsev/starship/inventory/internal/model"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository/converter"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) (map[string]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := converter.PartsFilterToRepo(filter)
	cursor, err := r.collection.Find(ctx, repoFilter)
	if err != nil {
		log.Printf("Ошибка получения данных по фильтру: фильтр - %v, ошибка - %v\n", filter, err)
	}
	defer func() {
		if cerr := cursor.Close(ctx); cerr != nil {
			log.Printf("Ошибка при закрытии курсора: %v\n", cerr)
		}
	}()

	var parts []model.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		log.Printf("Ошибка декодирования parts: %v\n", err)
	}
	list := map[string]model.Part{}
	for _, part := range parts {
		list[part.Uuid] = part
	}
	if len(list) == 0 {
		return map[string]model.Part{}, model.ErrPartListEmpty
	}
	return list, nil
}
