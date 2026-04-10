package order

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrLoadNotFound = errors.New("load not found")

type Repository interface {
	Create(ctx context.Context, load *Load) error
	ListByPosterID(ctx context.Context, posterID string) ([]Load, error)
	ListAll(ctx context.Context) ([]Load, error)
	ListUnassignedByPickupWindow(ctx context.Context, from, to time.Time) ([]Load, error)
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
			weight_kg, status, pickup_latitude, pickup_longitude, dropoff_latitude,
			dropoff_longitude, pickup_at, priority, cargo_type, equipment_type,
			quoted_price_kes, matched_at, matching_score, assignment_source,
			created_at, updated_at, picked_at, delivered_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20, $21,
			NOW(), NOW(), $22, $23
		)
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
		load.PickupLatitude,
		load.PickupLongitude,
		load.DropoffLatitude,
		load.DropoffLongitude,
		load.PickupAt,
		load.Priority,
		load.CargoType,
		load.EquipmentType,
		load.QuotedPriceKES,
		load.MatchedAt,
		load.MatchingScore,
		load.AssignmentSource,
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

func (r *PostgresRepository) ListUnassignedByPickupWindow(ctx context.Context, from, to time.Time) ([]Load, error) {
	rows, err := r.db.Query(ctx, baseLoadQuery()+`
		WHERE assigned_driver_id IS NULL
		  AND pickup_at >= $1
		  AND pickup_at < $2
		  AND status = $3
		ORDER BY priority DESC, pickup_at ASC, weight_kg DESC
	`, from, to, StatusPosted)
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
			pickup_latitude = $9,
			pickup_longitude = $10,
			dropoff_latitude = $11,
			dropoff_longitude = $12,
			pickup_at = $13,
			priority = $14,
			cargo_type = $15,
			equipment_type = $16,
			quoted_price_kes = $17,
			matched_at = $18,
			matching_score = $19,
			assignment_source = $20,
			updated_at = $21,
			picked_at = $22,
			delivered_at = $23
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
		load.PickupLatitude,
		load.PickupLongitude,
		load.DropoffLatitude,
		load.DropoffLongitude,
		load.PickupAt,
		load.Priority,
		load.CargoType,
		load.EquipmentType,
		load.QuotedPriceKES,
		load.MatchedAt,
		load.MatchingScore,
		load.AssignmentSource,
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
		&load.PickupLatitude,
		&load.PickupLongitude,
		&load.DropoffLatitude,
		&load.DropoffLongitude,
		&load.PickupAt,
		&load.Priority,
		&load.CargoType,
		&load.EquipmentType,
		&load.QuotedPriceKES,
		&load.MatchedAt,
		&load.MatchingScore,
		&load.AssignmentSource,
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
		       weight_kg, status, pickup_latitude, pickup_longitude, dropoff_latitude,
		       dropoff_longitude, pickup_at, priority, cargo_type, equipment_type,
		       quoted_price_kes, matched_at, matching_score, assignment_source,
		       created_at, updated_at, picked_at, delivered_at
		FROM loads
	`
}
