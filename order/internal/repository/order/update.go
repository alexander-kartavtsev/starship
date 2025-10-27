package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
)

func (r *repository) Update(ctx context.Context, uuid string, updateInfo model.OrderUpdateInfo) error {
	log.Println("Проверяем наличие заказа...")
	_, err := r.Get(ctx, uuid)
	if err != nil {
		return model.ErrOrderNotFound
	}
	log.Println("...заказ найден")

	r.mu.Lock()
	defer r.mu.Unlock()

	log.Println("Генерим запрос...")
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", updateInfo.Status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"order_uuid": uuid})

	if updateInfo.TransactionUuid != nil {
		builderUpdate = builderUpdate.Set("transaction_uuid", updateInfo.TransactionUuid)
	}
	if updateInfo.PaymentMethod != nil {
		builderUpdate = builderUpdate.Set("payment_method", updateInfo.PaymentMethod)
	}

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("Ошибка генерации запроса update: %v\n", err)
		return err
	}
	log.Println("...готово")

	log.Println("Выполняем запрос...")
	res, err := r.poolDb.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Ошибка обновления данных в таблице orders: %v\n", err)
		return err
	}
	log.Println("...готово")

	log.Printf("Обновлено количество строк: %d", res.RowsAffected())

	return nil
}
