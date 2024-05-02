CREATE TABLE IF NOT EXISTS shop_items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    in_stock INTEGER NOT NULL
);