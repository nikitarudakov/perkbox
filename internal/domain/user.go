package domain

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id"`
	BusinessID uuid.UUID `json:"business_id"` // A foreign key
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	Email      string    `json:"email"`
}
