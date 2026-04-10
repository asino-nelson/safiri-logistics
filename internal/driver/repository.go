package driver

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrProfileNotFound = errors.New("driver kyc profile not found")

type Repository interface {
	Upsert(ctx context.Context, profile *Profile) error
	GetByUserID(ctx context.Context, userID string) (*Profile, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Upsert(ctx context.Context, profile *Profile) error {
	query := `
		INSERT INTO driver_profiles (
			user_id, national_id, license_number, truck_registration, status,
			rejection_reason, submitted_at, reviewed_at, verified_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id)
		DO UPDATE SET national_id = EXCLUDED.national_id,
		              license_number = EXCLUDED.license_number,
		              truck_registration = EXCLUDED.truck_registration,
		              status = EXCLUDED.status,
		              rejection_reason = EXCLUDED.rejection_reason,
		              submitted_at = EXCLUDED.submitted_at,
		              reviewed_at = EXCLUDED.reviewed_at,
		              verified_at = EXCLUDED.verified_at
	`

	_, err := r.db.Exec(ctx, query,
		profile.UserID,
		profile.NationalID,
		profile.LicenseNumber,
		profile.TruckRegistration,
		profile.Status,
		profile.RejectionReason,
		profile.SubmittedAt,
		profile.ReviewedAt,
		profile.VerifiedAt,
	)
	return err
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*Profile, error) {
	query := `
		SELECT user_id, national_id, license_number, truck_registration, status,
		       rejection_reason, submitted_at, reviewed_at, verified_at
		FROM driver_profiles
		WHERE user_id = $1
	`

	var profile Profile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.NationalID,
		&profile.LicenseNumber,
		&profile.TruckRegistration,
		&profile.Status,
		&profile.RejectionReason,
		&profile.SubmittedAt,
		&profile.ReviewedAt,
		&profile.VerifiedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}

		return nil, err
	}

	return &profile, nil
}
