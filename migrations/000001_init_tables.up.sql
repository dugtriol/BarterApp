BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TYPE IF EXISTS user_mode CASCADE;
CREATE TYPE user_mode AS ENUM (
    'CLIENT',
    'ADMIN'
    );

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name     VARCHAR(50)        NOT NULL,
    email    VARCHAR(50) UNIQUE NOT NULL,
    phone    VARCHAR(50),
    password VARCHAR(100)       NOT NULL,
    city     VARCHAR(50),
    mode     user_mode
);

DROP TYPE IF EXISTS product_category CASCADE;
CREATE TYPE product_category AS ENUM (
    'HOME',
    'CLOTHES',
    'CHILDREN',
    'SPORT',
    'OTHER'
    );

DROP TYPE IF EXISTS product_status CASCADE;
CREATE TYPE product_status AS ENUM (
    'AVAILABLE',
    'EXCHANGING',
    'EXCHANGED'
    );

CREATE TABLE IF NOT EXISTS products
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(50) NOT NULL,
    description VARCHAR(50) NOT NULL,
    image       VARCHAR(150),
    status      product_status   DEFAULT 'AVAILABLE',
    category    product_category,
    user_id     UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at  TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS favorites
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    product_id UUID REFERENCES products (id) ON DELETE CASCADE
);

DROP TYPE IF EXISTS transaction_status CASCADE;
CREATE TYPE transaction_status AS ENUM (
    'CREATED',
    'ONGOING',
    'DONE',
    'DECLINED'
    );

DROP TYPE IF EXISTS shipping_method CASCADE;
CREATE TYPE shipping_method AS ENUM (
    'MEETUP',
    'MAIL',
    'COURIER'
    );

CREATE TABLE IF NOT EXISTS transactions
(
    id               UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    owner            UUID REFERENCES users (id) ON DELETE CASCADE,
    buyer            UUID REFERENCES users (id) ON DELETE CASCADE,
    product_id_owner UUID REFERENCES products (id) ON DELETE CASCADE,
    product_id_buyer UUID REFERENCES products (id) ON DELETE CASCADE,
    created_at       TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    shipping         shipping_method    DEFAULT 'MEETUP',
    address          VARCHAR(200) NOT NULL,
    status           transaction_status DEFAULT 'CREATED'
);

-- comments
-- conversation?
-- message

COMMIT;