package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/kafka"
	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/logger"
	eventsV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/events/v1"
)

type service struct {
	orderAssembledProducer kafka.Producer
}

func NewService(producer kafka.Producer) *service {
	return &service{orderAssembledProducer: producer}
}

func (s *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error {
	msg := &eventsV1.OrderAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembledMsg", zap.Error(err))
		return err
	}

	err = s.orderAssembledProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembledMsg", zap.Error(err))
		return err
	}

	return nil
}
