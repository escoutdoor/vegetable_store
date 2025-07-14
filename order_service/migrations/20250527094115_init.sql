-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS orders
(
    id uuid primary key default gen_random_uuid(),
    user_id uuid,
    total_amount decimal(10, 2)
);

CREATE TABLE IF NOT EXISTS recipients
(
    id uuid primary key default gen_random_uuid(),
    first_name text not null,
    last_name text not null,
    phone_number text not null,
    email text not null
);

CREATE TABLE IF NOT EXISTS addresses
(
    id uuid primary key default gen_random_uuid(),
    address text not null
);

CREATE TABLE IF NOT EXISTS order_items
(
    id uuid primary key default gen_random_uuid(),
    order_id uuid references orders(id),
    vegetable_id uuid,
    weight decimal(8, 3),
    price decimal(10, 2),
    discounted_price decimal(10, 2),
    recipient_id uuid references recipients(id),
    address_id uuid references addresses(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS recipients;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS orders;

-- +goose StatementEnd
