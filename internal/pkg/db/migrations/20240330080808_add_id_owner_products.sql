-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN id_owner UUID
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT products;
-- +goose StatementEnd
