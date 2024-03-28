-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN city VARCHAR(50) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT users;
-- +goose StatementEnd
