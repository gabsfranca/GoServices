package domain

type Invoice struct {
	Id       int64         `json:"id"`
	Nf       string        `json:"nf"`
	Status   string        `json:"status"`
	Type     string        `json:"type"`
	Products []InvoiceItem `json:"products"`
}

type InvoiceItem struct {
	Id           int64   `json:"id"`
	InvoiceId    int64   `json:"invoiceId"`
	SerialNumber string  `json:"serialNumber"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	TotalPrice   float64 `json:"totalPrice"`
}
