package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/converter"
	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

func (r *repository) Read(ctx context.Context, orderUuid string) (model.Order, error) {
	var order repomodel.Order
	orderUUID := uuid.MustParse(orderUuid)
	builderSelect := sq.
		Select(
			repomodel.ColumnOrderUUID,
			repomodel.ColumnUserUUID,
			repomodel.ColumnPartUuids,
			repomodel.ColumnTotalPrice,
			repomodel.ColumnTransactionUUID,
			repomodel.ColumnPaymentMethod,
			repomodel.ColumnStatus).
		From(repomodel.TableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{repomodel.ColumnOrderUUID: orderUUID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.Order{}, err
	}

	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		&order.PartUuids,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		log.Printf("failed to select order: %v\n", err)
		return model.Order{}, err
	}

	return converter.ModelToOrder(order), nil
}
