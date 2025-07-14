package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikitarudakov/perkbox/internal/domain"
	"log"
	"net/http"
)

type Repository interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id uuid.UUID) error
	GetUserByID(id uuid.UUID) (*domain.User, error)
	ListUsersForBusiness(businessID uuid.UUID) ([]domain.User, error)
}

type UserHandler struct {
	r Repository
}

func NewUserHandler(r Repository) *UserHandler {
	return &UserHandler{r: r}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	businessID, err := uuid.Parse(c.GetHeader("X-User-Business"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-User-Business header"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if businessID != user.BusinessID {
		c.JSON(http.StatusForbidden, gin.H{"message": "users must belong to the same business"})
		return
	}

	if err := h.r.CreateUser(&user); err != nil {
		log.Println("error creating user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id path parameter"})
		return
	}

	businessID, err := uuid.Parse(c.GetHeader("X-User-Business"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid X-User-Business header"})
		return
	}

	user, err := h.r.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	if user.BusinessID != businessID {
		c.JSON(http.StatusForbidden, gin.H{"message": "cannot delete user from a different business"})
		return
	}

	if err := h.r.DeleteUser(userID); err != nil {
		log.Println("error deleting user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user was deleted"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.GetHeader("X-User-Id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id in the header"})
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if user.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{"message": "users id mismatch"})
		return
	}

	if err := h.r.UpdateUser(&user); err != nil {
		log.Println("error listing all users: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id path parameter"})
		return
	}

	user, err := h.r.GetUserByID(userID)
	if err != nil {
		log.Println("error getting user by id: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access to admin details is restricted"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	businessIDParam := c.Param("business_id")
	headerBusinessID := c.GetHeader("X-User-Business")

	if businessIDParam != headerBusinessID {
		c.JSON(http.StatusForbidden, gin.H{"message": "business ID mismatch"})
		return
	}

	businessID, err := uuid.Parse(businessIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business id"})
		return
	}

	users, err := h.r.ListUsersForBusiness(businessID)
	if err != nil {
		log.Println("error listing users:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, users)
}
