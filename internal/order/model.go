package order

import "time"

const (
	StatusPosted    = "posted"
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
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	PickedAt         *time.Time `json:"picked_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
}

type CreateLoadRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	WeightKG    float64 `json:"weight_kg" binding:"required,gt=0"`
}

type UpdateLoadStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
