package converter

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	payment_v1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/payment/v1"
)

// ModelPayOrderToProtoPayOrder конвертирует заказ на оплату из модели в прото-модель.
func ModelPayOrderToProtoPayOrder(payOrder model.PayOrder) *payment_v1.PayOrderRequest {
	return &payment_v1.PayOrderRequest{
		OrderUuid:     payOrder.OrderUUID,
		UserUuid:      payOrder.UserUUID,
		PaymentMethod: modelPaymentMethodToProtoPaymentMethod(payOrder.PaymentMethod),
	}
}

// modelPaymentMethodToProtoPaymentMethod конвертирует метод оплаты из модели в прото-модель.
func modelPaymentMethodToProtoPaymentMethod(paymentMethod string) payment_v1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_SBP
	case model.PaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_INVESTOR_MONEY
	case model.PaymentMethodCARD:
		return payment_v1.PaymentMethod_CARD
	default:
		return payment_v1.PaymentMethod_UNKNOWN_UNSPECIFIED
	}
}
