-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    name       VARCHAR(50)  NOT NULL,
    lastname   VARCHAR(50)  NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(100)        NOT NULL,
    city       VARCHAR(100)        NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT users;
-- +goose StatementEnd