-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS vegetables
(
    id uuid primary key default gen_random_uuid(),
    name text,
    weight float
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS vegetables;

-- +goose StatementEnd
