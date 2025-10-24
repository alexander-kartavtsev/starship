package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderUuid := uuid.NewString()

	tx, err := r.poolDb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(err)
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Printf("Ошибка отмены tr: %v\n", err)
		}
	}()

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("order_uuid", "user_uuid", "total_price", "status", "payment_method").
		Values(orderUuid, order.UserUuid, order.TotalPrice, repoModel.PendingPayment, repoModel.Card).
		Suffix("RETURNING order_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("Ошибка build query: %v\n", err)
	}

	var orderUuidDb string
	err = r.poolDb.QueryRow(ctx, query, args...).Scan(&orderUuidDb)
	if err != nil {
		log.Printf("Ошибка insert в таблицу orders: %v\n", err)
	}

	builderInsert = sq.Insert("order_items").
		PlaceholderFormat(sq.Dollar).
		Columns("order_uuid", "part_uuid")

	for _, partUuid := range order.PartUuids {
		builderInsert = builderInsert.Values(orderUuidDb, partUuid)
	}

	query, args, err = builderInsert.ToSql()
	if err != nil {
		log.Printf("Ошибка build query: %v\n", err)
	}
	rows, err := r.poolDb.Query(ctx, query, args...)
	if err != nil {
		log.Printf("Ошибка insert в таблицу orders: %v\n", err)
	}
	rows.Close()
	log.Printf("Добавлена запись : %v\n", rows)

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}

	return orderUuid, nil
}
