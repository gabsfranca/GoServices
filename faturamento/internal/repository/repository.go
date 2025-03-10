package repository

import (
	"context"

	"github.com/gabsfranca/gerador-nf-faturamento/internal/domain"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, i *domain.Invoice) error
}
