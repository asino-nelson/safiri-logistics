package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrLoadNotFound = errors.New("load not found")

type Repository interface {
	Create(ctx context.Context, load *Load) error
	ListByPosterID(ctx context.Context, posterID string) ([]Load, error)
	ListAll(ctx context.Context) ([]Load, error)
	GetByID(ctx context.Context, loadID string) (*Load, error)
	Update(ctx context.Context, load *Load) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, load *Load) error {
	query := `
		INSERT INTO loads (
			id, poster_id, assigned_driver_id, title, description, origin, destination,
			weight_kg, status, created_at, updated_at, picked_at, delivered_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW(), $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		load.ID,
		load.PosterID,
		load.AssignedDriverID,
		load.Title,
		load.Description,
		load.Origin,
		load.Destination,
		load.WeightKG,
		load.Status,
		load.PickedAt,
		load.DeliveredAt,
	)
	return err
}

func (r *PostgresRepository) ListByPosterID(ctx context.Context, posterID string) ([]Load, error) {
	rows, err := r.db.Query(ctx, baseLoadQuery()+" WHERE poster_id = $1 ORDER BY created_at DESC", posterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanLoads(rows)
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]Load, error) {
	rows, err := r.db.Query(ctx, baseLoadQuery()+" ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanLoads(rows)
}

func (r *PostgresRepository) GetByID(ctx context.Context, loadID string) (*Load, error) {
	row := r.db.QueryRow(ctx, baseLoadQuery()+" WHERE id = $1", loadID)

	load, err := scanLoad(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrLoadNotFound
		}

		return nil, err
	}

	return load, nil
}

func (r *PostgresRepository) Update(ctx context.Context, load *Load) error {
	query := `
		UPDATE loads
		SET assigned_driver_id = $2,
			title = $3,
			description = $4,
			origin = $5,
			destination = $6,
			weight_kg = $7,
			status = $8,
			updated_at = $9,
			picked_at = $10,
			delivered_at = $11
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(ctx, query,
		load.ID,
		load.AssignedDriverID,
		load.Title,
		load.Description,
		load.Origin,
		load.Destination,
		load.WeightKG,
		load.Status,
		load.UpdatedAt,
		load.PickedAt,
		load.DeliveredAt,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrLoadNotFound
	}

	return nil
}

type loadScanner interface {
	Scan(dest ...any) error
}

func scanLoad(scanner loadScanner) (*Load, error) {
	var load Load
	err := scanner.Scan(
		&load.ID,
		&load.PosterID,
		&load.AssignedDriverID,
		&load.Title,
		&load.Description,
		&load.Origin,
		&load.Destination,
		&load.WeightKG,
		&load.Status,
		&load.CreatedAt,
		&load.UpdatedAt,
		&load.PickedAt,
		&load.DeliveredAt,
	)
	if err != nil {
		return nil, err
	}

	return &load, nil
}

func scanLoads(rows pgx.Rows) ([]Load, error) {
	loads := make([]Load, 0)
	for rows.Next() {
		load, err := scanLoad(rows)
		if err != nil {
			return nil, err
		}

		loads = append(loads, *load)
	}

	return loads, rows.Err()
}

func baseLoadQuery() string {
	return `
		SELECT id, poster_id, assigned_driver_id, title, description, origin, destination,
		       weight_kg, status, created_at, updated_at, picked_at, delivered_at
		FROM loads
	`
}
