CREATE DATABASE products_service_db;


\c products_service_db 

CREATE TABLE IF NOT EXISTS products(
    id SERIAL PRIMARY KEY,
    serial_number VARCHAR(30) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL, 
    current_stock INT NOT NULL DEFAULT 0,
);
