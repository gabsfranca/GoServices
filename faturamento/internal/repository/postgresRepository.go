package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gabsfranca/gerador-nf-faturamento/internal/domain"
)

type postgresInvoiceRepository struct {
	db *sql.DB
}

func NewPostgresInvoiceRepository(db *sql.DB) InvoiceRepository {
	return &postgresInvoiceRepository{db: db}
}

func (r *postgresInvoiceRepository) CreateInvoice(ctx context.Context, i *domain.Invoice) error {
	existingInvoice, _ := r.GetInvoiceByNF(ctx, i.Nf)
	if existingInvoice != nil {
		return errors.New("NF já cadastrada")
	}

	query := `
		INSERT INTO invoices (nf, status, type)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		i.Nf,
		i.Status,
		i.Type,
	).Scan(&i.Id)

	if err != nil {
		fmt.Println("eror na query: ", err)
		return err
	}

	itemQuery := `
		INSERT INTO invoice_items (
			invoice_id, 
			serial_number, 
			quantity,
			price, 
			discount,
			total_price
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, item := range i.Products {
		totalPrice := (item.Price * item.Quantity) * (1 - item.Discount/100)

		_, err = r.db.ExecContext(
			ctx,
			itemQuery,
			i.Id,
			item.SerialNumber,
			item.Quantity,
			item.Price,
			item.Discount,
			totalPrice,
		)

		if err != nil {
			return fmt.Errorf("erro ao inserir no item: ", err)
		}
	}

	return err

}

func (r *postgresInvoiceRepository) GetInvoiceById(ctx context.Context, id int64) (*domain.Invoice, error) {
	query := `
		SELECT id, nf, total_value, status 
		FROM invoices
		WHERE ID = $1
	`

	i := &domain.Invoice{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&i.Id,
		&i.Nf,
		&i.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("nf nao encontrada")
		}
		return nil, err
	}
	return i, nil

}

func (r *postgresInvoiceRepository) GetInvoiceByNF(ctx context.Context, NF string) (*domain.Invoice, error) {
	query := `
		SELECT id, nf, total_value, status, type
		FROM invoices
		WHERE nf = $1
	`

	i := &domain.Invoice{}
	err := r.db.QueryRowContext(ctx, query, NF).Scan(
		&i.Id,
		&i.Nf,
		&i.Status,
		&i.Type,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("nf nao encontrada")
		}
		return nil, err
	}
	return i, nil
}

func (r *postgresInvoiceRepository) UpdateInvoiceStatus(ctx context.Context, InvoiceId int64, status string) error {
	query := `
		UPDATE invoices
		SET status = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, InvoiceId)
	return err
}
