package repo

import (
	"L0/models"
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
)

type IDBRepo interface {
	InsertOrder(ctx context.Context, order *models.EventOrder) error
	GetOrderByID(ctx context.Context, id string) (*models.EventOrder, error)
	GetOrders(ctx context.Context) ([]models.EventOrder, error)
}

type dbRepo struct {
	db *sqlx.DB
}

// TODO fix SQL query to write into DB
func (r *dbRepo) InsertOrder(ctx context.Context, order *models.EventOrder) error {
	result, err := r.db.NamedQueryContext(ctx, `INSERT INTO orders (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) 
VALUES (:order_uid,:track_number,:entry,:locale,:internal_signature,:customer_id,:delivery_service,:shardkey,:sm_id,:date_created,:oof_shard) returning id`, order)
	if err != nil {
		return err
	}
	var orderID int64
	for result.Next() {
		if err := result.Scan(&orderID); err != nil {
			return err
		}
	}
	result.Close()

	_, err = r.db.ExecContext(ctx, `INSERT INTO deliveries (order_id,name,phone,zip,city,address,region,email) 
VALUES ($1, $2,$3,$4,$5,$6,$7,$8)`, orderID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	_, err = r.db.ExecContext(ctx, `INSERT INTO payments (order_id,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) 
VALUES ($1, $2,$3,$4,$5,$6,$7,$8, $9, $10, $11)`, orderID, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	result, err = r.db.NamedQueryContext(ctx, `INSERT INTO items(chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) values (:chrt_id,:track_number,:price,:rid,:name,:sale,:size,:total_price,:nm_id,:brand,:status) on conflict(chrt_id) do update set updated_at = now() returning id`, order.Items)
	if err != nil {
		return err
	}
	orderIDStr := strconv.FormatInt(orderID, 10)
	insertOrderItems := ""
	for result.Next() {
		var itemID int64
		if err := result.Scan(&itemID); err != nil {
			return err
		}

		insertOrderItems += "(" + orderIDStr + "," + strconv.FormatInt(itemID, 10) + "),"
	}
	insertOrderItems = insertOrderItems[:len(insertOrderItems)-1]
	result.Close()

	_, err = r.db.ExecContext(ctx, `INSERT INTO orders_items(order_id, item_id) values `+insertOrderItems)
	return err
}

func (r *dbRepo) GetOrderByID(ctx context.Context, id string) (*models.EventOrder, error) {
	row := r.db.QueryRowxContext(ctx, `SELECT json_build_object('id', o.id, 'order_uid', o.order_uid, 'track_number', track_number, 'entry', entry, 'locale', locale,
                         'internal_signature', internal_signature, 'customer_id', customer_id, 'delivery_service', delivery_service,
                         'shardkey', shardkey, 'sm_id', sm_id, 'date_created', date_created, 'oof_shard', oof_shard, 'created_at', o.created_at, 
                         'delivery', json_build_object('name', name, 'phone',phone, 'zip', zip, 'city', city, 'address', address, 'region', region, 'email', email, 'created_at', d.created_at ), 
                         'payment', json_build_object('transaction', transaction, 'request_id', request_id, 'currency', currency, 'provider', provider, 'amount', amount, 'payment_dt', payment_dt, 'bank', bank, 'delivery_cost', delivery_cost, 'goods_total', goods_total, 'custom_fee', custom_fee, 'created_at', p.created_at)) FROM orders o
                              left join deliveries d on d.order_id = o.id
                              left join payments p on p.order_id = o.id
         WHERE o.id = $1`, id)
	var order models.EventOrder
	var jsonData string
	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(jsonData), &order); err != nil {
		return nil, err
	}

	var items []models.Item
	rows, err := r.db.QueryxContext(ctx, `select it.* from items it left join orders_items oi on oi.item_id = it.id where order_id = $1`, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	log.Println(items)
	order.Items = items

	return &order, nil
}

func (r *dbRepo) GetOrders(ctx context.Context) ([]models.EventOrder, error) {
	rows, err := r.db.QueryxContext(ctx, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []models.EventOrder
	for rows.Next() {
		var order models.EventOrder
		err := rows.StructScan(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
