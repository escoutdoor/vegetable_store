-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id uuid primary key default gen_random_uuid(),
    first_name text,
    last_name text,
    email text unique,
    phone_number text,
    password text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
