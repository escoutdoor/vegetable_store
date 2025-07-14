-- +goose Up
-- +goose StatementBegin
ALTER TABLE vegetables ADD COLUMN price float;
ALTER TABLE vegetables ADD COLUMN discounted_price float;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vegetables DROP COLUMN price;
ALTER TABLE vegetables DROP COLUMN discounted_price;
-- +goose StatementEnd
