package domain

import "time"

type Invoice struct {
	Id         int64     `json:"id"`
	Nf         int64     `json:"nf"`
	IssueDate  time.Time `json:"issue_date"`
	TotalValue float64   `json:"totalValue"`
	Status     string    `json:"status"`
}

type InvoiceItem struct {
	Id           int64   `json:"id"`
	InvoiceId    int64   `json:"invoiceId"`
	SerialNumber string  `json:"serialNumer"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	TotalPrice   float64 `json:"totalPrice"`
	Discount     float64 `json:"discount"`
}
