package order

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alexander-kartavtsev/starship/order/internal/config"
)

func GetDbConn() *pgx.Conn {
	ctx := context.Background()

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, config.AppConfig().Postgres.Uri())
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		panic(err)
	}

	// Проверяем, что соединение с базой установлено
	err = con.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		panic(err)
	}

	return con
}

func GetDbPool() *pgxpool.Pool {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.Uri())
	if err != nil {
		log.Printf("Ошибка подключения к б/д: %v\n", err)
	}

	return pool
}
