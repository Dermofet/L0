DROP INDEX IF EXISTS order_uid_idx;
DROP INDEX IF EXISTS delivery_uid_idx;
DROP INDEX IF EXISTS rid_idx;
DROP INDEX IF EXISTS transaction_idx;

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;