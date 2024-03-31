-- +goose Up
-- +goose StatementBegin
CREATE TABLE products
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(50)  NOT NULL,
    description VARCHAR(500) NOT NULL,
    image       VARCHAR(500) NOT NULL,
    city        VARCHAR(100) NOT NULL,
    status      VARCHAR(100) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT products;
-- +goose StatementEnd
