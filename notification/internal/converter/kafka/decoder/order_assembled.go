package decoder

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/vipshark78/microservices-course-homeworks/notification/internal/model"
	eventsV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) DecodeOrderAssembled(data []byte) (model.OrderAssembledEvent, error) {
	var pb eventsV1.OrderAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse UserUuid: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse EventUuid: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderAssembledEvent{}, fmt.Errorf("failed to parse OrderUuid: %w", err)
	}

	return model.OrderAssembledEvent{
		EventUUID:    eventUUID.String(),
		OrderUUID:    orderUUID.String(),
		UserUUID:     userUUID.String(),
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
