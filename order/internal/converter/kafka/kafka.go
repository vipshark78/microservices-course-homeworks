package kafka

import "github.com/vipshark78/microservices-course-homeworks/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledEvent, error)
}
