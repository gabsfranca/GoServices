package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gabsfranca/gerador-nf-estoque/internal/domain"
)

type postgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *domain.Product) error {

	existingProduct, _ := r.GetBySerialNumber(ctx, product.SerialNumber)
	if existingProduct != nil {
		return errors.New("produto já cadastrado")
	}

	query := `
		INSERT INTO products (serial_number, name, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.SerialNumber,
		product.Name,
		product.Description,
		product.Price,
	).Scan(&product.ID, &product.CreatedAt)

	return err
}

func (r *postgresProductRepository) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, serial_number, name, description, price, current_stock, created_at
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.SerialNumber,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CurrentStock,
		&product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) GetBySerialNumber(ctx context.Context, serialNumber string) (*domain.Product, error) {
	query := `
		SELECT id, serial_number, name, description, price, current_stock, created_at
		FROM products
		WHERE serial_number = $1
	`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, serialNumber).Scan(
		&product.ID,
		&product.SerialNumber,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CurrentStock,
		&product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) GetProduts(ctx context.Context) ([]*domain.Product, error) {
	query := `
		SELECT id, serial_number, name, description, price, current_stock, created_at
		FROM products
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.SerialNumber,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CurrentStock,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err != nil {
		return nil, err
	}

	return products, nil

}
