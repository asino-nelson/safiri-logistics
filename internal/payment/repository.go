package payment

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, payment *Payment) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, payment *Payment) error {
	query := `
		INSERT INTO payments (
			id, load_id, customer_id, provider, provider_reference, phone_number,
			amount_kes, currency, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		payment.ID,
		payment.LoadID,
		payment.CustomerID,
		payment.Provider,
		payment.ProviderReference,
		payment.PhoneNumber,
		payment.AmountKES,
		payment.Currency,
		payment.Status,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	return err
}
