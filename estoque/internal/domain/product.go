package domain

import (
	"time"
)

type Product struct {
	ID           int64     `json:"id"`
	SerialNumber string    `json:"serialNumber"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	CurrentStock int       `json:"currentStock"`
	CreatedAt    time.Time `json:"created_at"`
}

type StockMovements struct {
	ID           int64     `json:"id"`
	ProductId    int64     `json:"product_id"`
	Quantity     int64     `json:"quantity"`
	MovementType string    `json:"movement_type"`
	InvoiceId    string    `json:"invoice_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
