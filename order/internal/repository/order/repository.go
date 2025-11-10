package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	repomodel "github.com/vipshark78/microservices-course-homeworks/order/internal/repository/model"
)

type repository struct {
	orders map[string]repomodel.Order
	pool   *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		orders: make(map[string]repomodel.Order),
		pool:   pool,
	}
}
