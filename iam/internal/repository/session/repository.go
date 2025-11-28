package session

import "github.com/alexander-kartavtsev/starship/platform/pkg/cache"

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}
