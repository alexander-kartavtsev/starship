package assembly_consumer

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/assembly/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

func (s *service) AssemblyHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.assemblyDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode Order in assembly service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("Uuid", event.Uuid),
		zap.String("OrderUuid", event.OrderUuid),
		zap.String("UserUuid", event.UserUuid),
		zap.String("PaymentMethod", string(event.PaymentMethod)),
		zap.String("TransactionUuid", event.TransactionUuid),
		zap.String("Type", event.Type),
	)

	start := time.Now()
	log.Println("Получили оплату. Начинаем процесс. Осталось...")
	for i := 10; i > 0; i-- {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			<-ctx.Done()
		}()
		log.Println("...", i-1, " секунд")
	}
	log.Println("Поехали!!!")
	buildTime := time.Since(start)

	err = s.assemblyProducerService.ProduceShipAssembled(ctx, model.ShipAssembledKafkaEvent{
		EventUuid:    event.Uuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: int64(buildTime.Seconds()),
	})
	if err != nil {
		return err
	}

	return nil
}
