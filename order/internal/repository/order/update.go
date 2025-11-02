package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/vipshark78/microservices-course-homeworks/order/internal/model"
	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	builderUpdate := sq.Update(repomodel.TableName).
		PlaceholderFormat(sq.Dollar).
		Set(repomodel.ColumnPartUuids, order.PartUuids).
		Set(repomodel.ColumnTotalPrice, order.TotalPrice).
		Set(repomodel.ColumnUpdatedAt, time.Now()).
		Set(repomodel.ColumnStatus, order.Status).
		Set(repomodel.ColumnTransactionUUID, order.TransactionUUID).
		Set(repomodel.ColumnPaymentMethod, order.PaymentMethod).
		Where(sq.Eq{repomodel.ColumnOrderUUID: order.OrderUUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.ErrInternalError
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update note: %v\n", err)
		return model.ErrInternalError
	}

	if res.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
