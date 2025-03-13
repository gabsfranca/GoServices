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

func (r *postgresProductRepository) Create(ctx context.Context, p *domain.Product) error {

	existingProduct, _ := r.GetProductBySerialNumber(ctx, p.SerialNumber)
	if existingProduct != nil {
		return errors.New("produto já cadastrado")
	}

	query := `
		INSERT INTO products (serial_number, name, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		p.SerialNumber,
		p.Name,
		p.Description,
		p.Price,
	).Scan(&p.ID)

	return err
}

func (r *postgresProductRepository) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, serial_number, name, description, price, current_stock
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
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) GetProductBySerialNumber(ctx context.Context, serialNumber string) (*domain.Product, error) {
	query := `
		SELECT id, serial_number, name, description, price, current_stock
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
		SELECT id, serial_number, name, description, price, current_stock
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

func (r *postgresProductRepository) UpdateStock(ctx context.Context, id int64, newStock int) error {
	query := `
		UPDATE products
		SET current_stock = $1
		WHERE id = $2;
	`

	_, err := r.db.ExecContext(ctx, query, newStock, id)
	return err
}
