package domain

import (
	"time"
)

type Product struct {
	ID           int64     `json:"id"`
	SerialNumber string    `json:"codigo"`
	Name         string    `json:"nome"`
	Desc         string    `json:"desc"`
	Price        float64   `json:"preco"`
	CurrentStock int       `json:"estoque_atual"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StockMovements struct {
	ID           int64     `json:"id"`
	ProductId    int64     `json:"id_product"`
	Quantity     int64     `json:"quantidade"`
	MovementType string    `json:"tipo_transacao"`
	Invoice_id   string    `json:"id_fatura,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
