-- Создание таблицы OrderDB
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    delivery_uid VARCHAR(255),
    payment_transaction VARCHAR(255),
    locale VARCHAR(255),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(255),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS order_uid_idx ON orders (order_uid);

-- Создание таблицы Item
CREATE TABLE IF NOT EXISTS items (
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    sale INT,
    size VARCHAR(255),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT
);

CREATE INDEX IF NOT EXISTS rid_idx ON items (rid);

-- Создание таблицы Payment
CREATE TABLE IF NOT EXISTS payments (
    transaction VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount INT,
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE INDEX IF NOT EXISTS transaction_idx ON payments (transaction);

-- Создание таблицы DeliveryDB
CREATE TABLE IF NOT EXISTS deliveries (
    delivery_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS delivery_uid_idx ON deliveries (delivery_uid);
