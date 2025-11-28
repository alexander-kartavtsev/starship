package user

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	connDb *pgx.Conn
	poolDb *pgxpool.Pool
}

func NewRepository(con *pgx.Conn, pool *pgxpool.Pool) *repository {
	return &repository{
		connDb: con,
		poolDb: pool,
	}
}
