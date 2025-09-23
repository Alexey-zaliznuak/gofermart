-- users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,   -- в копейках
    withdraw BIGINT NOT NULL DEFAULT 0   -- в копейках
);

-- orders
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    accrual BIGINT,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- withdrawals
CREATE TABLE withdrawals (
    id SERIAL PRIMARY KEY,
    number TEXT NOT NULL UNIQUE,
    sum BIGINT,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
