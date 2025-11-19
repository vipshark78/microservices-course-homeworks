package order

import (
	"context"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	"github.com/vipshark78/microservices-course-homeworks/order/internal/repository/converter"
	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

func (r *repository) Insert(ctx context.Context, userUuid string, partUuids []string, price float64, status model.OrderStatus) (model.Order, error) {
	builderInsert := sq.
		Insert(repomodel.TableName).
		PlaceholderFormat(sq.Dollar).
		Columns(
			repomodel.ColumnUserUUID,
			repomodel.ColumnPartUuids,
			repomodel.ColumnTotalPrice,
			repomodel.ColumnStatus,
		).
		Values(
			userUuid,
			partUuids,
			price,
			status,
		).
		Suffix(fmt.Sprintf("RETURNING %s", repomodel.ColumnOrderUUID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.Order{}, err
	}

	var orderUUID uuid.UUID
	err = r.pool.QueryRow(ctx, query, args...).Scan(&orderUUID)
	if err != nil {
		log.Printf("failed to insert order: %v\n", err)
		return model.Order{}, err
	}

	newOrder := repomodel.Order{
		UserUUID:        userUuid,
		OrderUUID:       orderUUID.String(),
		PartUuids:       partUuids,
		TotalPrice:      price,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          repomodel.OrderStatusPENDINGPAYMENT,
	}
	return converter.ModelToOrder(newOrder), nil
}
