package storage

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Userer interface {
	Create(ctx context.Context, u *models.UserDTO) error
	Update(ctx context.Context, u *models.UserDTO) error
	GetByID(ctx context.Context, userID int) (*models.UserDTO, error)
	GetByUsername(ctx context.Context, username string) (*models.UserDTO, error)
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
func (s *UserStorage) Create(ctx context.Context, u *models.UserDTO) error {
	sql, args, err := sq.Insert("users").PlaceholderFormat(sq.Dollar).
		Columns("username", "password_hash").
		Values(u.Username, u.Hash).ToSql()

	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return err
}
func (s *UserStorage) CreateTokenGPT(ctx context.Context, token string) error {
	sql, args, err := sq.Insert("users").PlaceholderFormat(sq.Dollar).
		Columns("token_gpt").Values(token).ToSql()

	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// Update - обновление пользователя в БД
func (s *UserStorage) Update(ctx context.Context, u *models.UserDTO) error {
	sql, args, err := sq.Update("users").PlaceholderFormat(sq.Dollar).
		SetMap(map[string]interface{}{
			"username":      u.Username,
			"password_hash": u.Hash,
			"token_gpt":     u.TokenGPT,
		}).Where(sq.Eq{"id": u.ID}).ToSql()

	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

// GetByID - получение пользователя по IDUser из БД
func (s *UserStorage) GetByID(ctx context.Context, userID int) (*models.UserDTO, error) {
	sql, args, err := sq.Select("id", "username", "password_hash", "token_gpt").PlaceholderFormat(sq.Dollar).
		From("users").Where(sq.Eq{"id": userID}).ToSql()
	if err != nil {
		return nil, err
	}
	u := &models.UserDTO{}
	err = s.conn.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Username, &u.Hash, &u.TokenGPT)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (*models.UserDTO, error) {
	sql, args, err := sq.Select("id", "username", "password_hash", "token_gpt").PlaceholderFormat(sq.Dollar).
		From("users").Where(sq.Eq{"username": username}).ToSql()
	if err != nil {
		return nil, err
	}
	u := &models.UserDTO{}
	_ = s.conn.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Username, &u.Hash, &u.TokenGPT)
	if err != nil {
		return nil, err
	}

	return u, nil

}
