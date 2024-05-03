CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT now(),
    hired_at TIMESTAMP NOT NULL DEFAULT now()
);
