package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/alexander-kartavtsev/starship/notification/internal/client/http"
	"github.com/alexander-kartavtsev/starship/notification/internal/model"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
)

const chatID = 1026089093

//go:embed templates/assembled_notification.tmpl
//go:embed templates/paid_notification.tmpl
var templateFS embed.FS

type assembledTemplateData struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type paidTemplateData struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   string
	TransactionUuid string
	Type            string
}

var (
	assembledTemplate = template.Must(template.ParseFS(templateFS, "templates/assembled_notification.tmpl"))
	paidTemplate      = template.Must(template.ParseFS(templateFS, "templates/paid_notification.tmpl"))
)

type service struct {
	telegramClient http.TelegramClient
}

// NewService создает новый Telegram сервис
func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendAssembledNotification(ctx context.Context, event model.ShipAssembledKafkaEvent) error {
	message, err := s.buildAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message (assembled) sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildUFOMessage создает сообщение о наблюдении UFO из шаблона
func (s *service) buildAssembledMessage(event model.ShipAssembledKafkaEvent) (string, error) {
	data := assembledTemplateData{
		EventUuid:    event.EventUuid,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendPaidNotification(ctx context.Context, event model.OrderKafkaEvent) error {
	message, err := s.buildPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message (assembled) sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildUFOMessage создает сообщение о наблюдении UFO из шаблона
func (s *service) buildPaidMessage(event model.OrderKafkaEvent) (string, error) {
	data := paidTemplateData{
		EventUuid:       event.Uuid,
		OrderUuid:       event.OrderUuid,
		UserUuid:        event.UserUuid,
		PaymentMethod:   string(event.PaymentMethod),
		TransactionUuid: event.TransactionUuid,
		Type:            event.Type,
	}

	var buf bytes.Buffer
	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
