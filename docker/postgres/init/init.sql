drop table IF EXISTS users_info;

create table orders
(
    id                 serial PRIMARY KEY,
    order_uid          text,
    track_number       text,
    entry              text,
    delivery_id        int,
    payment_id         int,
    locale             text default 'en',
    internal_signature text default '',
    customer_id        text default 'test',
    delivery_service   text default 'meest',
    shardkey           text default '0',
    sm_id              int  default 0,
    date_created       TIMESTAMP,
    oof_shard          text default '1'
);

create table delivery
(
    id      serial PRIMARY KEY,
    name    text,
    phone   text,
    zip     text,
    city    text,
    address text,
    region  text,
    email   text
);

create table payment
(
    id            serial PRIMARY KEY,
    transaction   text,
    request_id    text,
    currency      text,
    provider      text,
    amount        int default 0,
    payment_dt    int,
    bank          text,
    delivery_cost int default 0,
    goods_total   int default 0,
    custom_fee    int default 0
);

create table items
(
    id           serial PRIMARY KEY,
    order_id     int,
    chrt_id      int,
    track_number text,
    price        int,
    rid          text,
    name         text,
    sale         int,
    size         text,
    total_price  int,
    nm_id        int,
    brand        text,
    status       int
);