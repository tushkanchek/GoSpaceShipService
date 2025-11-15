-- +goose Up


CREATE TABLE orders(
    order_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID NOT NULL,
    part_uuids  UUID[] NOT NULL,
    total_price DOUBLE PRECISION,
    transaction_uuid UUID,
	payment_method TEXT,
	order_status TEXT NOT NULL

);

-- +goose Down
DROP TABLE IF EXISTS orders;
