package storage

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/errors"
	"github.com/KuYaki/waffler_server/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Userer interface {
	Create(ctx context.Context, u *models.User) (int, error)
	Update(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetByLogin(ctx context.Context, username string) (models.User, error)
}

// UserStorage - хранилище пользователей
type UserStorage struct {
	conn *pgxpool.Pool
}

func NewUserStorage(conn *pgxpool.Pool) Userer {
	return &UserStorage{conn: conn}
}

const (
	userCacheKey     = "user:%d"
	userCacheTTL     = 15
	userCacheTimeout = 50
)

// Create - создание пользователя в БД
func (s *UserStorage) Create(ctx context.Context, u *models.User) (int, error) {
	sql, args, err := sq.Insert("users").Columns("username", "hash_pass").Values(u.Username, u.Pass).ToSql()

	if err != nil {
		return errors.InternalError, err
	}
	_, err = s.conn.Query(ctx, sql, args...)
	if err != nil {
		return errors.InternalError, err
	}

	return 0, err
}

// Update - обновление пользователя в БД
func (s *UserStorage) Update(ctx context.Context, u *models.User) error {
	log.Fatal("Implemented Update in UserStorage")
	return nil
}

// GetByID - получение пользователя по ID из БД
func (s *UserStorage) GetByID(ctx context.Context, userID int) (*models.User, error) {
	log.Fatal("Implemented GetByID in UserStorage")
	return &models.User{}, nil
}

func (s *UserStorage) GetByLogin(ctx context.Context, username string) (models.User, error) {
	log.Fatal("Implemented GetByLogin in UserStorage")
	return models.User{}, nil

}
