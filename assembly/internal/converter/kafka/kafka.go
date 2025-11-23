package kafka

import "github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
