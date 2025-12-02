package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/vipshark78/microservices-course-homeworks/notification/internal/client/http"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/config"
	"github.com/vipshark78/microservices-course-homeworks/notification/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
)

//go:embed templates/paid_notification.tmpl templates/assembled_notification.tmpl
var templateFS embed.FS

type paidTemplateData struct {
	OrderUUID       string
	TransactionUUID string
	PaymentMethod   string
}

type assembledTemplateData struct {
	OrderUUID    string
	BuildTimeSec int64
}

var (
	paidTemplate      = template.Must(template.ParseFS(templateFS, "templates/paid_notification.tmpl"))
	assembledTemplate = template.Must(template.ParseFS(templateFS, "templates/assembled_notification.tmpl"))
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

// SendOrderPaidNotification отправляет уведомление об оплате заказа
func (s *service) SendOrderPaidNotification(ctx context.Context, notification model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(notification)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, config.AppConfig().TelegramBot.ChatID(), message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", config.AppConfig().TelegramBot.ChatID()), zap.String("message", message))
	return nil
}

// SendOrderAssembledNotification отправляет уведомление о сборке заказа
func (s *service) SendOrderAssembledNotification(ctx context.Context, notification model.OrderAssembledEvent) error {
	message, err := s.buildOrderAssembledMessage(notification)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, config.AppConfig().TelegramBot.ChatID(), message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", config.AppConfig().TelegramBot.ChatID()), zap.String("message", message))
	return nil
}

// buildOrderAssembledMessage создает сообщение о сборке заказа из шаблона
func (s *service) buildOrderAssembledMessage(notification model.OrderAssembledEvent) (string, error) {
	data := assembledTemplateData{
		OrderUUID:    notification.OrderUUID,
		BuildTimeSec: notification.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// buildOrderPaidMessage создает сообщение об оплате заказа из шаблона
func (s *service) buildOrderPaidMessage(notification model.OrderPaidEvent) (string, error) {
	data := paidTemplateData{
		OrderUUID:       notification.OrderUUID,
		TransactionUUID: notification.TransactionUUID,
		PaymentMethod:   notification.PaymentMethod,
	}

	var buf bytes.Buffer
	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
