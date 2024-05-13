CREATE TABLE IF NOT EXISTS transaction_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO transaction_types (name) VALUES
    ('transfer'),
    ('purchase'),
    ('deposit'),
    ('refund'),
    ('reward')
    ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS transaction_statuses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO transaction_statuses (name) VALUES
    ('pending'),
    ('completed'),
    ('cancelled')
    ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL REFERENCES users(id),
    receiver_id INTEGER NOT NULL REFERENCES users(id),
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    type_id INTEGER NOT NULL REFERENCES transaction_types(id),
    status_id INTEGER NOT NULL REFERENCES transaction_statuses(id)
)