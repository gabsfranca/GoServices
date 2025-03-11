CREATE DATABASE products_service_db;


\c products_service_db 

CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    serial_number VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL, 
    current_stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS stock_movements(
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    movement_type VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT stock_movements_product_fk
    FOREIGN KEY (product_id) REFERENCES products(id)
);