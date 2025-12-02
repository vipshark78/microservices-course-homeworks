package service

import (
	"context"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledEvent) error
}
