package order

import (
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/alexander-kartavtsev/starship/order/internal/migrator"
)

type repository struct {
	mu     sync.RWMutex
	connDb *pgx.Conn
	poolDb *pgxpool.Pool
}

func NewRepository(con *pgx.Conn, pool *pgxpool.Pool) *repository {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationsDir)
	err := migratorRunner.Up()
	// err := migratorRunner.Down()
	// err = migratorRunner.Down()
	if err != nil {
		log.Printf("Ошибка миграции базы данных: %v\n", err)
		panic(err)
	}

	return &repository{
		connDb: con,
		poolDb: pool,
	}
}
