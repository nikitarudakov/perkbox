package repo

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nikitarudakov/perkbox/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host string `env:"POSTGRES_HOST"`
	Port string `env:"POSTGRES_PORT"`
	User string `env:"POSTGRES_USER"`
	Pass string `env:"POSTGRES_PASSWORD"`
	DB   string `env:"USER_DB"`
}

func LoadConfig() (*Config, error) {
	// Load .env file (optional, for dev only)
	_ = godotenv.Load(".env")

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(cfg *Config) (*Repository, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
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
	return r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *Repository) DeleteUser(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *Repository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := r.db.Model(&domain.User{}).
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) ListUsersForBusiness(businessID uuid.UUID) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.Model(&domain.User{}).
		Where("business_id = ? and role = 'user'", businessID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
