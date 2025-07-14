package repo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nikitarudakov/perkbox/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host string `env:"POSTGRES_HOST"`
	Port string `env:"POSTGRES_PORT"`
	User string `env:"POSTGRES_USER"`
	Pass string `env:"POSTGRES_PASS"`
	DB   string `env:"USER_DB"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *Config) (*Repository, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) CreateUser(user *domain.User) error {
	return r.db.Model(&domain.User{}).Create(user).Error
}

func (r *Repository) UpdateUser(user *domain.User) error {
	return r.db.Model(&domain.User{}).Updates(user).Error
}

func (r *Repository) DeleteUser(user *domain.User) error {
	return r.db.Model(&domain.User{}).Delete(user).Error
}

func (r *Repository) ListAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Model(&domain.User{}).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) ListUsersForBusiness(businessID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Model(&domain.User{}).Find(&users, "business_id = ?", businessID).Error; err != nil {
		return nil, err
	}
	return users, nil
}
