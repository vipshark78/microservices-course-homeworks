package v1

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

// CreateOrder реализует обработку запроса создания заказа.
func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	if err := a.validateCreateOrderRequest(req); err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	reqUuids, err := a.convertUUIDToSliceString(req.PartUuids)
	if err != nil {
		return &order_v1.BadRequestError{Code: 400, Message: "Bad Request"}, err
	}

	order, err := a.orderService.CreateOrder(ctx, req.UserUUID.String(), reqUuids)
	if err != nil {
		return &order_v1.InternalServerError{Code: 500, Message: "Internal Server Error"}, err
	}
	return &order_v1.CreateOrderResponse{OrderUUID: uuid.MustParse(order.OrderUUID), TotalPrice: order.TotalPrice}, nil
}

// validateCreateOrderRequest выполняет валидацию входящего запроса.
func (a *api) validateCreateOrderRequest(req *order_v1.CreateOrderRequest) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("ошибка валидации запроса: %w", err)
	}

	if err := uuid.Validate(req.UserUUID.String()); err != nil {
		return fmt.Errorf("ошибка валидации UUID пользователя: %w", err)
	}
	return nil
}

// convertUUIDToSliceString конвертирует массив UUID в строковый формат.
func (a *api) convertUUIDToSliceString(uuids []uuid.UUID) ([]string, error) {
	strUuids := make([]string, 0, len(uuids))
	for _, UUID := range uuids {
		stringUUID := UUID.String()
		if err := uuid.Validate(stringUUID); err != nil {
			return nil, fmt.Errorf("ошибка валидации UUID: %w", err)
		}
		strUuids = append(strUuids, stringUUID)
	}
	return strUuids, nil
}
