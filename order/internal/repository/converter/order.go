package converter

import (
	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

// ModelToOrder конвертирует модель репозитория в модель бизнес-логики.
func ModelToOrder(order repomodel.Order) model.Order {
	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		Status:          model.OrderStatus(order.Status),
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   ModelToPaymentMethod(order.PaymentMethod),
	}
}

func PaymentMethodToModel(paymentMethod *model.PaymentMethod) *repomodel.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}
	pm := repomodel.PaymentMethod(*paymentMethod)
	return &pm
}

func ModelToPaymentMethod(paymentMethod *repomodel.PaymentMethod) *model.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}
	pm := model.PaymentMethod(*paymentMethod)
	return &pm
}

// OrderToModel конвертирует модель бизнес-логики в модель репозитория.
func OrderToModel(order model.Order) repomodel.Order {
	return repomodel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		Status:          repomodel.OrderStatus(order.Status),
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToModel(order.PaymentMethod),
	}
}
