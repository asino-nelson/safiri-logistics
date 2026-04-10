package order

import "time"

const (
	StatusPosted    = "posted"
	StatusMatched   = "matched"
	StatusPicked    = "picked"
	StatusInTransit = "in_transit"
	StatusDelivered = "delivered"
)

type Load struct {
	ID               string     `json:"id"`
	PosterID         string     `json:"poster_id"`
	AssignedDriverID *string    `json:"assigned_driver_id,omitempty"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Origin           string     `json:"origin"`
	Destination      string     `json:"destination"`
	WeightKG         float64    `json:"weight_kg"`
	Status           string     `json:"status"`
	PickupLatitude   float64    `json:"pickup_latitude"`
	PickupLongitude  float64    `json:"pickup_longitude"`
	DropoffLatitude  float64    `json:"dropoff_latitude"`
	DropoffLongitude float64    `json:"dropoff_longitude"`
	PickupAt         time.Time  `json:"pickup_at"`
	Priority         int        `json:"priority"`
	CargoType        string     `json:"cargo_type"`
	EquipmentType    string     `json:"equipment_type"`
	QuotedPriceKES   float64    `json:"quoted_price_kes"`
	MatchedAt        *time.Time `json:"matched_at,omitempty"`
	MatchingScore    *float64   `json:"matching_score,omitempty"`
	AssignmentSource string     `json:"assignment_source"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	PickedAt         *time.Time `json:"picked_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
}

type CreateLoadRequest struct {
	Title            string     `json:"title" binding:"required"`
	Description      string     `json:"description" binding:"required"`
	Origin           string     `json:"origin" binding:"required"`
	Destination      string     `json:"destination" binding:"required"`
	WeightKG         float64    `json:"weight_kg" binding:"required,gt=0"`
	PickupLatitude   float64    `json:"pickup_latitude" binding:"required"`
	PickupLongitude  float64    `json:"pickup_longitude" binding:"required"`
	DropoffLatitude  float64    `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude float64    `json:"dropoff_longitude" binding:"required"`
	PickupAt         *time.Time `json:"pickup_at"`
	Priority         int        `json:"priority" binding:"omitempty,min=1,max=5"`
	CargoType        string     `json:"cargo_type"`
	EquipmentType    string     `json:"equipment_type"`
}

type UpdateLoadStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
