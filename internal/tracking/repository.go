package tracking

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, event *Event) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, event *Event) error {
	query := `
		INSERT INTO tracking_events (
			id, load_id, driver_id, latitude, longitude, speed_kph, heading_degrees, recorded_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		event.ID,
		event.LoadID,
		event.DriverID,
		event.Latitude,
		event.Longitude,
		event.SpeedKPH,
		event.HeadingDegrees,
		event.RecordedAt,
	)
	return err
}
