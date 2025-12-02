package decoder

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/vipshark78/microservices-course-homeworks/assembly/internal/model"
	eventsV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsV1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	userUUID, err := uuid.Parse(pb.UserUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse UserUuid: %w", err)
	}

	eventUUID, err := uuid.Parse(pb.EventUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse EventUuid: %w", err)
	}

	orderUUID, err := uuid.Parse(pb.OrderUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse OrderUuid: %w", err)
	}

	transactionUUID, err := uuid.Parse(pb.TransactionUuid)
	if err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to parse TransactionUuid: %w", err)
	}

	if pb.PaymentMethod == "" {
		return model.OrderPaidEvent{}, errors.New("bad PaymentMethod")
	}

	return model.OrderPaidEvent{
		EventUUID:       eventUUID.String(),
		OrderUUID:       orderUUID.String(),
		UserUUID:        userUUID.String(),
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: transactionUUID.String(),
	}, nil
}
