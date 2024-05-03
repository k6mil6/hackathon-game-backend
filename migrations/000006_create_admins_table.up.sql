CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO roles (name) VALUES
    ('admin')
    ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT DEFAULT '',
    password_hash TEXT NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT now(),
    registered_by INTEGER REFERENCES admins(id),
    role_id INTEGER NOT NULL REFERENCES roles(id)
);

INSERT INTO admins (username, password_hash, role_id) VALUES
    ('admin', '$2a$10$BaFkWXSyqxFj/yQP5lzxaeF4sEyoL8aYT85kYOo5ZyQuogMDxhsJG', 1)
    ON CONFLICT (username) DO NOTHING
-- log admin pass admin