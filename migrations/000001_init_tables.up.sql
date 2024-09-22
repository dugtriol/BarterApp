BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TYPE IF EXISTS user_mode CASCADE;
CREATE TYPE user_mode AS ENUM (
    'CLIENT',
    'ADMIN'
    );

CREATE TABLE IF NOT EXISTS users
(
    id    UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    phone VARCHAR(50),
    password VARCHAR(100) NOT NULL,
    city  VARCHAR(50),
    mode  user_mode
);

DROP TYPE IF EXISTS product_category CASCADE;
CREATE TYPE product_category AS ENUM (
    'HOME',
    'CLOTHES'
    );

DROP TYPE IF EXISTS product_status CASCADE;
CREATE TYPE product_status AS ENUM (
    'CREATED',
    'SOLD'
    );

CREATE TABLE IF NOT EXISTS products
(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name  VARCHAR(50) NOT NULL,
    description VARCHAR(50) NOT NULL,
    image  VARCHAR(150),
    status product_status,
    category product_category,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
-- comments
-- transaction
-- conversation?
-- message
-- favorites?
COMMIT;