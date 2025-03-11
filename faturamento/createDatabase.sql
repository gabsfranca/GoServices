CREATE DATABASE invoice_service_db;

\c invoice_service_db

CREATE TABLE IF NOT EXISTS invoices(
    id SERIAL PRIMARY KEY,
    nf VARCHAR(30) NOT NULL, 
    total_value DECIMAL(10,2),
    status VARCHAR(10) NOT NULL,
    type VARCHAR(5) NOT NULL
);

create TABLE IF NOT EXISTS invoice_items(
    id SERIAL PRIMARY KEY,
    invoice_id INT NOT NULL,
    serial_number VARCHAR(30) NOT NULL, 
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL, 
    total_price DECIMAL(10,2) NOT NULL,
    discount DECIMAL(10,2),

    CONSTRAINT invoice_items_invoices_fk
    FOREIGN KEY(invoice_id)  REFERENCES invoices(id)
)