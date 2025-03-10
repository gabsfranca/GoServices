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
		INSERT INTO produtos (codigo, nome, description, preco, estoque_atual)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		product.SerialNumber,
		product.Name,
		product.Description,
		product.Price,
		product.CurrentStock,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func (r *postgresProductRepository) GetProductById(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, codigo, nome, description, preco, estoque_atual, created_at, updated_at
		FROM produtos
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
		&product.UpdatedAt,
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
		SELECT id, codigo, nome, description, preco, estoque_atual, created_at, updated_at
		FROM produtos
		WHERE codigo = $1
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
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *postgresProductRepository) UpdateStock(
	ctx context.Context,
	id int64,
	quantity int,
	movementType string,
	invoiceID string,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var currentStock int
	stockQuery := "SELECT current_stock FROM produtos WHERE id = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, stockQuery, id).Scan(&currentStock)
	if err != nil {
		return err
	}

	var newStock int
	if movementType == "IN" {
		newStock = currentStock + quantity
	} else if movementType == "OUT" {
		if currentStock < quantity {
			return errors.New("ESTOQUE INSUFICIENTE")
		}

		newStock = currentStock - quantity
	} else {
		return errors.New("invalid movement type")
	}

	updateQuery := `
		UPDATE produtos
		SET current_stock = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = tx.ExecContext(ctx, updateQuery, newStock, id)
	if err != nil {
		return err
	}

	movementQuery := `
		INSERT INTO stock_movements (product_id, quantity, movement_type, invoice_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.ExecContext(ctx, movementQuery, id, quantity, movementType, invoiceID)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (r *postgresProductRepository) GetProduts(ctx context.Context) ([]*domain.Product, error) {
	query := `
		SELECT id, codigo, nome, description, preco, estoque_atual, created_at, updated_at
		FROM produtos
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
			&product.UpdatedAt,
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
