package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nikitarudakov/perkbox/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Mock Repo Implementation ---

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetUserByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if user := args.Get(0); user != nil {
		return user.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepo) CreateUser(u *domain.User) error                          { return nil }
func (m *MockRepo) UpdateUser(u *domain.User) error                          { return nil }
func (m *MockRepo) DeleteUser(id uuid.UUID) error                            { return nil }
func (m *MockRepo) ListUsersForBusiness(id uuid.UUID) ([]domain.User, error) { return nil, nil }

// --- Test Suite Definition ---

type UserHandlerTestSuite struct {
	suite.Suite
	router   *gin.Engine
	mockRepo *MockRepo
}

func (s *UserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.mockRepo = new(MockRepo)
	userHandler := NewUserHandler(s.mockRepo)

	r := gin.Default()
	r.GET("/api/users/:user_id", userHandler.GetUser)
	s.router = r
}

func (s *UserHandlerTestSuite) TestGetUser_Success() {
	id := uuid.New()
	user := &domain.User{
		ID:         id,
		BusinessID: uuid.New(),
		Name:       "Test User",
		Email:      "test@example.com",
		Role:       "user",
	}
	s.mockRepo.On("GetUserByID", id).Return(user, nil)

	req := httptest.NewRequest("GET", "/api/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestGetUser_AdminBlocked() {
	id := uuid.New()
	admin := &domain.User{
		ID:         id,
		BusinessID: uuid.New(),
		Name:       "Admin",
		Email:      "admin@example.com",
		Role:       "admin",
	}
	s.mockRepo.On("GetUserByID", id).Return(admin, nil)

	req := httptest.NewRequest("GET", "/api/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
}

func (s *UserHandlerTestSuite) TestGetUser_InvalidID() {
	req := httptest.NewRequest("GET", "/api/users/invalid-uuid", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *UserHandlerTestSuite) TestGetUser_RepoError() {
	id := uuid.New()
	s.mockRepo.On("GetUserByID", id).Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/users/"+id.String(), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
}

// --- Run the test suite ---

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
