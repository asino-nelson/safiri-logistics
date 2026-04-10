package payment

import "time"

const (
	StatusInitiated = "initiated"
	StatusPaid      = "paid"
	StatusFailed    = "failed"
)

type Payment struct {
	ID                string    `json:"id"`
	LoadID            string    `json:"load_id"`
	CustomerID        string    `json:"customer_id"`
	Provider          string    `json:"provider"`
	ProviderReference string    `json:"provider_reference"`
	PhoneNumber       string    `json:"phone_number"`
	AmountKES         float64   `json:"amount_kes"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CheckoutRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ProviderRequest struct {
	PhoneNumber string
	AmountKES   float64
	AccountRef  string
	Description string
}

type ProviderResponse struct {
	Provider          string
	ProviderReference string
	Status            string
}
