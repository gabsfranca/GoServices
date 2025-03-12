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
