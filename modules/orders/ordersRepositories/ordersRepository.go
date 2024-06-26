package ordersRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/Kamila3820/go-shop-tutorial/modules/orders"
	"github.com/Kamila3820/go-shop-tutorial/modules/orders/ordersPatterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",
			"o"."transfer_slip",
			"o"."status",
			(
				SELECT 
					array_to_json(array_agg("pt"))
				FROM (
					SELECT
						"spo"."id",
						"spo"."quantity",
						"spo"."products"
					FROM "product_order" "spo"
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",
			"o"."address",
			"o"."contact",
			(
				SELECT
					SUM(COALESCE(("po"."products"->>'price')::FLOAT*("po"."quantity")::FLOAT, 0))
				FROM "product_order" "po"
				WHERE "po"."order_id" = "o"."id"
			) AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "order" "o"
		WHERE "o"."id" = $1
	) AS "t";
	`

	orderData := &orders.Order{
		TransferSlip: &orders.Transferslip{},
		Products:     make([]*orders.ProductsOrder, 0),
	}

	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("get order failed: %v", err)
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return orderData, nil
}

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := ordersPatterns.FindOrderBuilder(r.db, req)
	engineer := ordersPatterns.FindOrderEngineer(builder)

	result := engineer.FindOrder()
	count := engineer.CountOrder()

	return result, count
}
