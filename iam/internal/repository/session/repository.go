package session

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alexander-kartavtsev/starship/platform/pkg/cache"
)

const (
	sessionKeyPrefix    = "session:"
	sessionSetKeyPrefix = "session:set:"
)

type repository struct {
	cache  cache.RedisClient
	poolDb *pgxpool.Pool
}

func NewRepository(cache cache.RedisClient, poolDb *pgxpool.Pool) *repository {
	return &repository{
		cache:  cache,
		poolDb: poolDb,
	}
}
