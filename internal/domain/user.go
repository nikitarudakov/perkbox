package domain

import "github.com/google/uuid"

type Role int

const (
	Default Role = iota
	Admin
)

type User struct {
	ID         uuid.UUID `json:"id"`
	BusinessID uuid.UUID `json:"business_id"` // A foreign key
	Name       string    `json:"name"`
	Role       Role      `json:"role"`
	Email      string    `json:"email"`
}
