package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kimvlry/simple-order-service/internal/domain"
	"github.com/kimvlry/simple-order-service/internal/interfaces"
)

type PgOrderRepo struct {
	db *sqlx.DB
}

func NewPgOrderRepo(db *sqlx.DB) interfaces.OrderRepository {
	return &PgOrderRepo{db: db}
}

func (r *PgOrderRepo) GetById(ctx context.Context, uid string) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE order_uid = $1", uid); err != nil {
		return nil, err
	}
	if err := r.db.GetContext(ctx, &order.Delivery, "SELECT * FROM deliveries WHERE order_uid = $1", uid); err != nil {
		return nil, err
	}
	if err := r.db.SelectContext(ctx, &order.Items, "SELECT * FROM items WHERE order_uid = $1", uid); err != nil {
		return nil, err
	}
	if err := r.db.GetContext(ctx, &order.Payment, "SELECT * FROM payments WHERE order_uid = $1", uid); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PgOrderRepo) GetAll(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	rows, err := r.db.QueryxContext(ctx, `SELECT * FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o domain.Order
		if err := rows.StructScan(&o); err != nil {
			continue
		}

		_ = r.db.GetContext(ctx, &o.Delivery, `SELECT * FROM deliveries WHERE order_uid = $1`, o.OrderUid)
		_ = r.db.GetContext(ctx, &o.Payment, `SELECT * FROM payments WHERE order_uid = $1`, o.OrderUid)
		_ = r.db.SelectContext(ctx, &o.Items, `SELECT * FROM items WHERE order_uid = $1`, o.OrderUid)

		orders = append(orders, o)
	}
	return orders, nil
}

func (r *PgOrderRepo) Save(ctx context.Context, order *domain.Order) (err error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err = r.insertOrder(tx, order); err != nil {
		return err
	}
	if err = r.insertDelivery(ctx, tx, &order.Delivery, order.OrderUid); err != nil {
		return err
	}
	if err = r.insertPayment(ctx, tx, &order.Payment, order.OrderUid); err != nil {
		return err
	}
	if err = r.replaceItems(ctx, tx, order.Items, order.OrderUid); err != nil {
		return err
	}
	return nil
}

func (r *PgOrderRepo) insertOrder(tx *sqlx.Tx, order *domain.Order) error {
	query := `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES (
            :order_uid, :track_number, :entry, :locale, :internal_signature,
            :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard
        )
        ON CONFLICT (order_uid) DO NOTHING
        RETURNING order_uid;
    `
	rows, err := tx.NamedQuery(query, order)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&order.OrderUid)
	}
	return nil
}

func (r *PgOrderRepo) insertDelivery(ctx context.Context, tx *sqlx.Tx, d *domain.Delivery, orderUid string) error {
	data := struct {
		OrderUID string `db:"order_uid"`
		*domain.Delivery
	}{
		OrderUID: orderUid,
		Delivery: d,
	}

	_, err := tx.NamedExecContext(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES (
            :order_uid, :name, :phone, :zip, :city, :address, :region, :email
        )
        ON CONFLICT (order_uid) DO UPDATE SET
            name = EXCLUDED.name,
            phone = EXCLUDED.phone,
            zip = EXCLUDED.zip,
            city = EXCLUDED.city,
            address = EXCLUDED.address,
            region = EXCLUDED.region,
            email = EXCLUDED.email;
    `, data)
	return err
}

func (r *PgOrderRepo) insertPayment(ctx context.Context, tx *sqlx.Tx, p *domain.Payment, orderUid string) error {
	data := struct {
		OrderUID string `db:"order_uid"`
		*domain.Payment
	}{
		OrderUID: orderUid,
		Payment:  p,
	}

	_, err := tx.NamedExecContext(ctx, `
        INSERT INTO payments (
            order_uid, transaction, request_id, currency, provider,
            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES (
            :order_uid, :transaction, :request_id, :currency, :provider,
            :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee
        )
        ON CONFLICT (order_uid) DO UPDATE SET
            transaction = EXCLUDED.transaction,
            request_id = EXCLUDED.request_id,
            currency = EXCLUDED.currency,
            provider = EXCLUDED.provider,
            amount = EXCLUDED.amount,
            payment_dt = EXCLUDED.payment_dt,
            bank = EXCLUDED.bank,
            delivery_cost = EXCLUDED.delivery_cost,
            goods_total = EXCLUDED.goods_total,
            custom_fee = EXCLUDED.custom_fee;
    `, data)
	return err
}

func (r *PgOrderRepo) replaceItems(ctx context.Context, tx *sqlx.Tx, items []domain.Item, orderUid string) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM items WHERE order_uid = $1`, orderUid)
	if err != nil {
		return err
	}

	for _, item := range items {
		data := struct {
			OrderUID string `db:"order_uid"`
			domain.Item
		}{
			OrderUID: orderUid,
			Item:     item,
		}

		_, err := tx.NamedExecContext(ctx, `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid,
                name, sale, size, total_price, nm_id, brand, status
            ) VALUES (
                :order_uid, :chrt_id, :track_number, :price, :rid,
                :name, :sale, :size, :total_price, :nm_id, :brand, :status
            )
        `, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PgOrderRepo) Close() error {
	return r.db.Close()
}
