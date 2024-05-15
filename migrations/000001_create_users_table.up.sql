CREATE TABLE IF NOT EXISTS classes (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO classes (name) VALUES
    ('cat'),
    ('dog'),
    ('racoon');


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    class_id INTEGER REFERENCES classes(id),
    registered_at TIMESTAMP NOT NULL DEFAULT now(),
    hired_at TIMESTAMP NOT NULL DEFAULT now()
);
