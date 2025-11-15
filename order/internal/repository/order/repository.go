package order

import (
	"github.com/jackc/pgx/v5"
	def "order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *pgx.Conn
}

func NewOrderRepository(db *pgx.Conn) *repository {
	return &repository{
		db: db,
	}
}
