package repository

import (
	"context"

	"github.com/gabsfranca/gerador-nf-faturamento/internal/domain"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, i *domain.Invoice) error
	GetInvoiceByNF(ctx context.Context, NG string) (*domain.Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID int64, status string) error
}
