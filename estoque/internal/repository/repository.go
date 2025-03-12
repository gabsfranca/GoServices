package repository

import (
	"context"

	"github.com/gabsfranca/gerador-nf-estoque/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetProductById(ctx context.Context, id int64) (*domain.Product, error)
	GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Product, error)
	GetProduts(ctx context.Context) ([]*domain.Product, error)
	UpdateStock(ctx context.Context, id int64, newStock int) error
}
