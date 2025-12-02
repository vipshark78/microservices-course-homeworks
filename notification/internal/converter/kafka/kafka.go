package kafka

import "github.com/vipshark78/microservices-course-homeworks/notification/internal/model"

type OrderAssembledDecoder interface {
	DecodeOrderAssembled(data []byte) (model.OrderAssembledEvent, error)
}

type OrderPaidDecoder interface {
	DecodeOrderPaid(data []byte) (model.OrderPaidEvent, error)
}
