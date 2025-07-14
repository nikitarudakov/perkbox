package handlers

import (
	"github.com/google/uuid"
	"github.com/nikitarudakov/perkbox/internal/domain"
)

type Repository interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(user *domain.User) error
	ListAllUsers() ([]domain.User, error)
	ListUsersForBusiness(businessID uuid.UUID) ([]domain.User, error)
}

type UserHandler struct {
	r Repository
}

func NewUserHandler(r Repository) *UserHandler {
	return &UserHandler{r: r}
}
