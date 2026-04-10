package user

import "time"

const (
	RoleCustomer = "customer"
	RoleDriver   = "driver"
	RoleAdmin    = "admin"
)

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
