package order

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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

	order, err := getOrder(ctx, r, orderUuid)
	if err != nil {
		return model.Order{}, err
	}

	partUuids, err := getPartUuids(ctx, r, orderUuid)
	if err != nil {
		log.Printf("Ошибка при получении списка parts: %v\n", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		panic(err)
	}

	order.PartUuids = partUuids
	return order, nil
}

func getOrder(ctx context.Context, r *repository, uuid string) (model.Order, error) {
	builderQuery := sq.Select(
		"order_uuid",
		"user_uuid",
		"total_price",
		"status",
		"transaction_uuid",
		"payment_method",
	).
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_uuid": uuid}).
		Limit(1)

	query, args, err := builderQuery.ToSql()
	if err != nil {
		log.Printf("Ошибка генерации запроса SELECT (Get): %v\n", err)
		return model.Order{}, model.ErrOrderNotFound
	}

	var orderUuid, userUuid, transactionUuid string
	var totalPrice float64
	var paymentMethod model.PaymentMethod
	var status model.OrderStatus

	err = r.poolDb.
		QueryRow(ctx, query, args...).
		Scan(
			&orderUuid,
			&userUuid,
			&totalPrice,
			&status,
			&transactionUuid,
			&paymentMethod,
		)
	if err != nil {
		log.Printf("Ошибка получения данных из таблицы orders: %v\n", err)
		return model.Order{}, model.ErrOrderNotFound
	}

	return model.Order{
		OrderUuid:       orderUuid,
		UserUuid:        userUuid,
		TotalPrice:      totalPrice,
		TransactionUuid: transactionUuid,
		PaymentMethod:   paymentMethod,
		Status:          status,
	}, nil
}

func getPartUuids(ctx context.Context, r *repository, orderUuid string) ([]string, error) {
	sqlSelect := "part_uuid"
	tblOrders := "order_items"
	sql := fmt.Sprintf("select %s from %s where order_uuid = $1", sqlSelect, tblOrders)

	rows, err := r.connDb.Query(ctx, sql, orderUuid)
	if err != nil {
		log.Printf("Ошибка получения данных из б/д: %v\n", err)
		return []string{}, err
	}
	defer rows.Close()

	var partUuids []string
	for rows.Next() {
		var partUuid string
		err = rows.Scan(&partUuid)
		if err != nil {
			log.Printf("Ошибка сканирования данных строки order_items: %v\n", err)
		}
		partUuids = append(partUuids, partUuid)
	}
	return partUuids, nil
}
