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
	ListAvailableApproved(ctx context.Context) ([]Profile, error)
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
			rejection_reason, submitted_at, reviewed_at, verified_at,
			years_experience, max_load_kg, equipment_type, is_online, is_available,
			current_latitude, current_longitude, last_location_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (user_id)
		DO UPDATE SET national_id = EXCLUDED.national_id,
		              license_number = EXCLUDED.license_number,
		              truck_registration = EXCLUDED.truck_registration,
		              status = EXCLUDED.status,
		              rejection_reason = EXCLUDED.rejection_reason,
		              submitted_at = EXCLUDED.submitted_at,
		              reviewed_at = EXCLUDED.reviewed_at,
		              verified_at = EXCLUDED.verified_at,
		              years_experience = EXCLUDED.years_experience,
		              max_load_kg = EXCLUDED.max_load_kg,
		              equipment_type = EXCLUDED.equipment_type,
		              is_online = EXCLUDED.is_online,
		              is_available = EXCLUDED.is_available,
		              current_latitude = EXCLUDED.current_latitude,
		              current_longitude = EXCLUDED.current_longitude,
		              last_location_at = EXCLUDED.last_location_at
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
		profile.YearsExperience,
		profile.MaxLoadKG,
		profile.EquipmentType,
		profile.IsOnline,
		profile.IsAvailable,
		profile.CurrentLatitude,
		profile.CurrentLongitude,
		profile.LastLocationAt,
	)
	return err
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*Profile, error) {
	query := `
		SELECT user_id, national_id, license_number, truck_registration, status,
		       rejection_reason, submitted_at, reviewed_at, verified_at,
		       years_experience, max_load_kg, equipment_type, is_online, is_available,
		       current_latitude, current_longitude, last_location_at
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
		&profile.YearsExperience,
		&profile.MaxLoadKG,
		&profile.EquipmentType,
		&profile.IsOnline,
		&profile.IsAvailable,
		&profile.CurrentLatitude,
		&profile.CurrentLongitude,
		&profile.LastLocationAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}

		return nil, err
	}

	return &profile, nil
}

func (r *PostgresRepository) ListAvailableApproved(ctx context.Context) ([]Profile, error) {
	query := `
		SELECT user_id, national_id, license_number, truck_registration, status,
		       rejection_reason, submitted_at, reviewed_at, verified_at,
		       years_experience, max_load_kg, equipment_type, is_online, is_available,
		       current_latitude, current_longitude, last_location_at
		FROM driver_profiles
		WHERE status = $1
		  AND is_online = TRUE
		  AND is_available = TRUE
		  AND current_latitude IS NOT NULL
		  AND current_longitude IS NOT NULL
		ORDER BY years_experience DESC, last_location_at DESC NULLS LAST
	`

	rows, err := r.db.Query(ctx, query, KYCStatusApproved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drivers := make([]Profile, 0)
	for rows.Next() {
		var profile Profile
		if err := rows.Scan(
			&profile.UserID,
			&profile.NationalID,
			&profile.LicenseNumber,
			&profile.TruckRegistration,
			&profile.Status,
			&profile.RejectionReason,
			&profile.SubmittedAt,
			&profile.ReviewedAt,
			&profile.VerifiedAt,
			&profile.YearsExperience,
			&profile.MaxLoadKG,
			&profile.EquipmentType,
			&profile.IsOnline,
			&profile.IsAvailable,
			&profile.CurrentLatitude,
			&profile.CurrentLongitude,
			&profile.LastLocationAt,
		); err != nil {
			return nil, err
		}

		drivers = append(drivers, profile)
	}

	return drivers, rows.Err()
}
