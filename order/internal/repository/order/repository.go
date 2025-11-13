package order

import (

	def "order/internal/repository"

	"github.com/jackc/pgx/v5"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *pgx.Conn
	
}

func NewOrderRepository(db *pgx.Conn) *repository{
	return &repository{
		db: db,
	}
}
