package domain

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID `json:"id" gorm:"default:gen_random_uuid()"`
	BusinessID uuid.UUID `json:"business_id"` // A foreign key
	Name       string    `json:"name"`
	Role       string    `json:"role" gorm:"default:'user'"`
	Email      string    `json:"email"`
}
