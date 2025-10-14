package order

import (
	"context"

	"github.com/alexander-kartavtsev/starship/order/internal/model"
	repoModel "github.com/alexander-kartavtsev/starship/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, uuid string, updateInfo model.OrderUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	if updateInfo.PaymentMethod != nil && *updateInfo.PaymentMethod != model.Unknown {
		order.PaymentMethod = repoModel.PaymentMethod(*updateInfo.PaymentMethod)
	}

	if updateInfo.TransactionUuid != nil {
		order.TransactionUuid = *updateInfo.TransactionUuid
	}

	if updateInfo.PartUuids != nil {
		for _, partUuid := range *updateInfo.PartUuids {
			order.PartUuids = append(order.PartUuids, partUuid)
		}
	}

	if updateInfo.Status != nil {
		newStatus := repoModel.OrderStatus(*updateInfo.Status)
		if newStatus == repoModel.Cancelled && order.Status == repoModel.Paid {
			return model.ErrCancelPaidOrder
		}
		order.Status = repoModel.OrderStatus(*updateInfo.Status)
	}

	r.data[uuid] = order

	return nil
}
