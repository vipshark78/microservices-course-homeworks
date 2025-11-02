package converter

import (
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	order_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/openapi/order/v1"
)

// ConvertModelToOrder конвертирует модель бизнес-логики в DTO для клиента.
func ConvertModelToOrder(model model.Order) order_v1.OrderDto {
	return order_v1.OrderDto{
		OrderUUID:       uuid.MustParse(model.OrderUUID),
		UserUUID:        uuid.MustParse(model.UserUUID),
		PartUuids:       ConvertModelToPartsUUIDs(model.PartUuids),
		TotalPrice:      model.TotalPrice,
		Status:          ConvertModelToOrderStatus(model.Status),
		PaymentMethod:   ConvertModelToPaymentMethod(model.PaymentMethod),
		TransactionUUID: ConvertModelToOptNilUUID(model.TransactionUUID),
	}
}

// ConvertModelToPartsUUIDs конвертирует partUUIDs из модели бизнес-логики в модель DTO.
func ConvertModelToPartsUUIDs(modelPartsUUIDs []string) []uuid.UUID {
	var partsUUIDs []uuid.UUID
	for _, modelPartUUID := range modelPartsUUIDs {
		partUUID := uuid.MustParse(modelPartUUID)
		partsUUIDs = append(partsUUIDs, partUUID)
	}
	return partsUUIDs
}

// ConvertModelToOrderStatus конвертирует статус заказа из модели бизнес-логики в модель DTO.
func ConvertModelToOrderStatus(modelStatus model.OrderStatus) order_v1.OrderStatus {
	switch modelStatus {
	case model.OrderStatusPENDINGPAYMENT:
		return order_v1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return order_v1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return order_v1.OrderStatusCANCELLED
	}
	return order_v1.OrderStatusPENDINGPAYMENT
}

// ConvertModelToPaymentMethod конвертирует метод оплаты из модели бизнес-логики в модель DTO.
func ConvertModelToPaymentMethod(modelPaymentMethod *model.PaymentMethod) order_v1.OptNilPaymentMethod {
	if modelPaymentMethod == nil {
		return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodUNKNOWN)
	}

	switch *modelPaymentMethod {
	case model.PaymentMethodSBP:
		return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodSBP)
	case model.PaymentMethodCARD:
		return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodCARD)
	case model.PaymentMethodCREDITCARD:
		return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodCREDITCARD)
	case model.PaymentMethodINVESTORMONEY:
		return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodINVESTORMONEY)
	}
	return order_v1.NewOptNilPaymentMethod(order_v1.PaymentMethodUNKNOWN)
}

// ConvertModelToOptNilUUID конвертирует transactionUUID из модели DTO в модель бизнес-логики.
func ConvertPaymentMethodToModel(paymentMethod order_v1.NilPaymentMethod) model.PaymentMethod {
	if paymentMethod.Null {
		return model.PaymentMethodUNKNOWN
	}
	switch paymentMethod.Value {
	case order_v1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case order_v1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case order_v1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case order_v1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

// ConvertModelToOptNilUUID конвертирует transactionUUID из модели бизнес-логики в модель DTO.
func ConvertModelToOptNilUUID(modelTransactionUUID *string) order_v1.OptNilUUID {
	if modelTransactionUUID == nil {
		return order_v1.NewOptNilUUID(uuid.Nil)
	} else {
		return order_v1.NewOptNilUUID(uuid.MustParse(*modelTransactionUUID))
	}
}
