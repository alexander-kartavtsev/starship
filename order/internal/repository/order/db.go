package order

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const envPath = "../deploy/compose/order/.env"

func GetDbConn() *pgx.Conn {
	ctx := context.Background()

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		panic(err)
	}

	dbURI := os.Getenv("DB_URI")

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbURI)
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

	dbURI := os.Getenv("DB_URI")

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("Ошибка подключения к б/д: %v\n", err)
	}

	return pool
}
