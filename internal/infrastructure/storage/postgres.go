package storage

import (
	"L0/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

const (
	maxConn        = 10
	acquireTimeout = time.Minute
)

type PgsStorage struct {
	conn *pgx.ConnPool
}

func DSN() string {
	return "postgresql://wb:wb@localhost:5432/wb"
}

func NewPgsStorage(ctx context.Context) (*PgsStorage, error) {
	conf, err := pgx.ParseURI(DSN())
	if err != nil {
		return nil, err
	}

	poolConf := &pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: maxConn,
		AcquireTimeout: acquireTimeout,
	}
	conn, err := pgx.NewConnPool(*poolConf)
	if err != nil {
		return nil, err
	}

	return &PgsStorage{conn: conn}, nil
}

const queryReceiveOrder string = `select 
       ord.order_uid, ord.track_number, ord.entry, ord.locale, ord.internal_signature, 
       ord.customer_id, ord.delivery_service, ord.shardkey, ord.sm_id, ord.date_created, ord.oof_shard,
       
       del.name, del.phone, del.zip, del.city, del.address, del.region, del.email, 
       
       pay.transaction, pay.request_id, pay.currency, pay.provider, pay.amount, pay.payment_dt, pay.bank, 
       pay.delivery_cost, pay.goods_total, pay.custom_fee
	
	from orders ord 
    left join delivery as del on ord.delivery_id = del.id 
    left join payment as pay on ord.payment_id = pay.id
	where ord.order_uid = $1`

const queryReceiveAllOrders string = `select 
       ord.id, ord.order_uid, ord.track_number, ord.entry, ord.locale, ord.internal_signature, 
       ord.customer_id, ord.delivery_service, ord.shardkey, ord.sm_id, ord.date_created, ord.oof_shard,
       
       del.name, del.phone, del.zip, del.city, del.address, del.region, del.email, 
       
       pay.transaction, pay.request_id, pay.currency, pay.provider, pay.amount, pay.payment_dt, pay.bank, 
       pay.delivery_cost, pay.goods_total, pay.custom_fee
	
	from orders ord 
    left join delivery as del on ord.delivery_id = del.id 
    left join payment as pay on ord.payment_id = pay.id`

const queryReceiveItems string = `select 
		it.chrt_id, it.track_number, it.price, it.rid, it.name, it.sale, it.size, it.total_price, it.nm_id, it.brand, it.status from items it where order_id = $1`

func (s *PgsStorage) ReceiveOrder(orderUID string) (domain.Order, error) {
	//implement me
	var out domain.Order

	err := s.conn.QueryRow(queryReceiveOrder, orderUID).Scan(&out.OrderUid, &out.TrackNumber, &out.Entry, &out.Locale, &out.InternalSignature,
		&out.CustomerId, &out.DeliveryService, &out.Shardkey, &out.SmId, &out.DateCreated, &out.OofShard,

		&out.Delivery.Name, &out.Delivery.Phone, &out.Delivery.Zip, &out.Delivery.City, &out.Delivery.Address, &out.Delivery.Region, &out.Delivery.Email,

		&out.Payment.Transaction, &out.Payment.RequestId, &out.Payment.Currency, &out.Payment.Provider, &out.Payment.Amount, &out.Payment.PaymentDt, &out.Payment.Bank,
		&out.Payment.DeliveryCost, &out.Payment.GoodsTotal, &out.Payment.CustomFee)
	if err != nil {
		return domain.Order{}, err
	}

	rows, err := s.conn.Query(queryReceiveItems, orderUID)
	if err != nil {
		return domain.Order{}, err
	}
	for rows.Next() {
		var item domain.Item
		_ = rows.Scan(
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)
		out.Items = append(out.Items, item)
	}
	return out, nil
}

func (s *PgsStorage) ReceiveAll() ([]domain.Order, error) {
	var orders []domain.Order
	var orderID int64

	rows, err := s.conn.Query(queryReceiveAllOrders)
	if err != nil {
		return []domain.Order{}, err
	}

	for rows.Next() {
		var out domain.Order

		err = rows.Scan(&orderID, &out.OrderUid, &out.TrackNumber, &out.Entry, &out.Locale, &out.InternalSignature,
			&out.CustomerId, &out.DeliveryService, &out.Shardkey, &out.SmId, &out.DateCreated, &out.OofShard,

			&out.Delivery.Name, &out.Delivery.Phone, &out.Delivery.Zip, &out.Delivery.City, &out.Delivery.Address, &out.Delivery.Region, &out.Delivery.Email,

			&out.Payment.Transaction, &out.Payment.RequestId, &out.Payment.Currency, &out.Payment.Provider, &out.Payment.Amount, &out.Payment.PaymentDt, &out.Payment.Bank,
			&out.Payment.DeliveryCost, &out.Payment.GoodsTotal, &out.Payment.CustomFee)
		if err != nil {
			return []domain.Order{}, err
		}

		rows, err := s.conn.Query(queryReceiveItems, orderID)
		if err != nil {
			return []domain.Order{}, err
		}
		for rows.Next() {
			var item domain.Item
			_ = rows.Scan(
				&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
				&item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
			)
			out.Items = append(out.Items, item)
		}
		orders = append(orders, out)
	}
	return orders, nil

}

const insertPayment string = `
	insert into payment(
        transaction, request_id, currency, provider, amount, payment_dt, 
	                    bank, delivery_cost, goods_total, custom_fee
    ) values
          ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
          returning id`

const insertDelivery string = `
	insert into delivery(
        name, phone, zip, city, address, region, 
	                    email
    ) values
          ($1, $2, $3, $4, $5, $6, $7)
          returning id`

const insertOrder string = `
	insert into orders(
        order_uid, track_number, entry, delivery_id, payment_id, locale, 
	                    internal_signature, customer_id, delivery_service, 
	                    shardkey, sm_id, date_created, oof_shard
    ) values
          ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
          returning id`

const insertItems string = `
	insert into items(
        order_id, chrt_id, track_number, price, rid, name, 
	                    sale, size, total_price, nm_id, brand, status
    ) values
          ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

func (s *PgsStorage) Insert(in domain.Order) error {
	tx, err := s.conn.Begin()
	if err != nil {
		fmt.Printf("error while create tx %s", err.Error())
	}

	var paymentID int64
	err = tx.QueryRow(insertPayment, in.Payment.Transaction, in.Payment.RequestId, in.Payment.Currency, in.Payment.Provider, in.Payment.Amount,
		in.Payment.PaymentDt, in.Payment.Bank, in.Payment.DeliveryCost, in.Payment.GoodsTotal, in.Payment.CustomFee).Scan(&paymentID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while tx.QueryRow(insertPayment) %w", err)
	}

	var deliveryID int64
	err = tx.QueryRow(insertDelivery, in.Delivery.Name, in.Delivery.Phone, in.Delivery.Zip, in.Delivery.City,
		in.Delivery.Address, in.Delivery.Region, in.Delivery.Email).Scan(&deliveryID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while tx.QueryRow(insertDelivery) %w", err)
	}

	var orderID int64
	err = tx.QueryRow(insertOrder, in.OrderUid, in.TrackNumber, in.Entry, deliveryID, paymentID, in.Locale, in.InternalSignature,
		in.CustomerId, in.DeliveryService, in.Shardkey, in.SmId, in.DateCreated, in.OofShard).Scan(&orderID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while tx.QueryRow(insertOrder) %w", err)
	}

	for _, item := range in.Items {
		_, err = tx.Exec(insertItems, orderID, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error while tx.QueryRow(insertItems): %w; on item: %+v", err, item)
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error while tx.Commit: %w", err)
	}
	return nil
}
