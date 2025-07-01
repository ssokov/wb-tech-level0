CREATE TABLE orders
(
    order_uid          TEXT PRIMARY KEY,
    track_number       TEXT NOT NULL,
    entry              TEXT NOT NULL,
    locale             TEXT NOT NULL,
    internal_signature TEXT NOT NULL,
    customer_id        TEXT NOT NULL,
    delivery_service   TEXT NOT NULL,
    shardkey           TEXT,
    sm_id              BIGINT NOT NULL,
    date_created       TIMESTAMPTZ NOT NULL,
    oof_shard          TEXT
);

CREATE TABLE deliveries
(
    order_uid TEXT PRIMARY KEY REFERENCES orders (order_uid) ON DELETE CASCADE,
    name      TEXT NOT NULL,
    phone     TEXT NOT NULL,
    zip       TEXT NOT NULL,
    city      TEXT NOT NULL,
    address   TEXT NOT NULL,
    region    TEXT NOT NULL,
    email     TEXT NOT NULL
);

CREATE TABLE payments
(
    order_uid     TEXT PRIMARY KEY REFERENCES orders (order_uid) ON DELETE CASCADE,
    transaction   TEXT NOT NULL,
    request_id    TEXT NOT NULL,
    currency      TEXT NOT NULL,
    provider      TEXT NOT NULL,
    amount        BIGINT NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          TEXT NOT NULL,
    delivery_cost BIGINT NOT NULL,
    goods_total   BIGINT NOT NULL,
    custom_fee    BIGINT NOT NULL
);

CREATE TABLE items
(
    id           BIGSERIAL PRIMARY KEY,
    order_uid    TEXT NOT NULL REFERENCES orders (order_uid) ON DELETE CASCADE,
    chrt_id      BIGINT NOT NULL,
    track_number TEXT NOT NULL,
    price        BIGINT NOT NULL,
    rid          TEXT NOT NULL,
    name         TEXT NOT NULL,
    sale         BIGINT NOT NULL,
    size         TEXT NOT NULL,
    total_price  BIGINT NOT NULL,
    nm_id        BIGINT NOT NULL,
    brand        TEXT NOT NULL,
    status       BIGINT NOT NULL
);
