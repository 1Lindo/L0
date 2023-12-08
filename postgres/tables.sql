create table orders
(
    id                 serial PRIMARY KEY,
    order_uid          varchar,
    track_number       varchar,
    entry              varchar,
    locale             varchar,
    internal_signature varchar,
    customer_id        varchar,
    delivery_service   varchar,
    shardkey           varchar,
    sm_id              int,
    date_created       timestamp,
    oof_shard          varchar,
    created_at         timestamp default now()
);

create table payments
(
    id            serial PRIMARY KEY,
    order_id      int,
    transaction   varchar,
    request_id    varchar,
    currency      varchar,
    provider      varchar,
    amount        int,
    payment_dt    int,
    bank          varchar,
    delivery_cost int,
    goods_total   int,
    custom_fee    int,
    created_at    timestamp default now()
);

create table deliveries
(
    id         serial PRIMARY KEY,
    order_id   int,
    name       varchar,
    phone      varchar,
    zip        varchar,
    city       varchar,
    address    varchar,
    region     varchar,
    email      varchar,
    created_at timestamp default now()
);
create table items
(
    id          serial PRIMARY KEY,
    chrt_id     int unique,
    track_number varchar,
    price       int,
    rid         varchar,
    name        varchar,
    sale        varchar,
    size        varchar,
    total_price int,
    nm_id       int,
    brand       varchar,
    status      int,
    updated_at  timestamp default now(),
    created_at  timestamp default now()
);

create table orders_items
(
    id       SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    item_id  INTEGER NOT NULL
);