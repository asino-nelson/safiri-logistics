package driver

import "time"

const (
	KYCStatusPending  = "pending"
	KYCStatusApproved = "approved"
	KYCStatusRejected = "rejected"
)

type Profile struct {
	UserID            string     `json:"user_id"`
	NationalID        string     `json:"national_id"`
	LicenseNumber     string     `json:"license_number"`
	TruckRegistration string     `json:"truck_registration"`
	Status            string     `json:"status"`
	RejectionReason   *string    `json:"rejection_reason,omitempty"`
	SubmittedAt       time.Time  `json:"submitted_at"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty"`
	YearsExperience   int        `json:"years_experience"`
	MaxLoadKG         float64    `json:"max_load_kg"`
	EquipmentType     string     `json:"equipment_type"`
	IsOnline          bool       `json:"is_online"`
	IsAvailable       bool       `json:"is_available"`
	CurrentLatitude   *float64   `json:"current_latitude,omitempty"`
	CurrentLongitude  *float64   `json:"current_longitude,omitempty"`
	LastLocationAt    *time.Time `json:"last_location_at,omitempty"`
}

type SubmitKYCRequest struct {
	NationalID        string `json:"national_id" binding:"required"`
	LicenseNumber     string `json:"license_number" binding:"required"`
	TruckRegistration string `json:"truck_registration" binding:"required"`
}

type ReviewKYCRequest struct {
	Status          string `json:"status" binding:"required"`
	RejectionReason string `json:"rejection_reason"`
}

type UpdateOperationsRequest struct {
	YearsExperience int      `json:"years_experience" binding:"gte=0"`
	MaxLoadKG       float64  `json:"max_load_kg" binding:"gte=0"`
	EquipmentType   string   `json:"equipment_type" binding:"required"`
	IsOnline        bool     `json:"is_online"`
	IsAvailable     bool     `json:"is_available"`
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
}
