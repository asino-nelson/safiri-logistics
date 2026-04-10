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
