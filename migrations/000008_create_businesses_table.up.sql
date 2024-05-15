CREATE TABLE businesses_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    profit DECIMAL(10, 2)
);

INSERT INTO businesses_types (name, description, profit) VALUES
    ('farm', 'Grows food', 10.00),
    ('factory', 'Produces food', 15.00),
    ('warehouse', 'Stores food', 20.00),
    ('bank', 'Stores money', 30.00);


CREATE TABLE businesses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type_id INTEGER REFERENCES businesses_types(id),
    price DECIMAL(10, 2),
    owner_id INTEGER REFERENCES users(id)
);

INSERT INTO businesses (name, type_id, price) VALUES
    ('Farm', 1, 1000.00),
    ('Factory', 2, 2000.00),
    ('Warehouse', 3, 2500.00),
    ('Bank', 4, 4000.00);

