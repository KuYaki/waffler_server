package storage

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/models"
	"gorm.io/gorm"
)

type Userer interface {
	Create(ctx context.Context, u *models.UserDTO) error
	Update(ctx context.Context, u *models.UserDTO) error
	GetByID(ctx context.Context, userID int) (*models.UserDTO, error)
	GetByUsername(ctx context.Context, username string) (*models.UserDTO, error)
}

// UserStorage - хранилище пользователей
type UserStorage struct {
	conn *gorm.DB
}

func NewUserStorage(conn *gorm.DB) Userer {
	return &UserStorage{conn: conn}
}

const (
	userCacheKey     = "user:%d"
	userCacheTTL     = 15
	userCacheTimeout = 50
)

// Create - создание пользователя в БД
func (s *UserStorage) Create(ctx context.Context, u *models.UserDTO) error {
	err := s.conn.Create(u).Error
	if err != nil {
		return err
	}

	return err
}

// Update - обновление пользователя в БД
func (s *UserStorage) Update(ctx context.Context, u *models.UserDTO) error {
	err := s.conn.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}

// GetByID - получение пользователя по IDUser из БД
func (s *UserStorage) GetByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	u := &models.UserDTO{}
	err := s.conn.Take(u, userID).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*models.UserDTO, error) {
	u := &models.UserDTO{}
	err := s.conn.Where("username = ?", username).Take(u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}
