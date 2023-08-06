package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/errors"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type Userer interface {
	Create(ctx context.Context, user models.User) int
	Update(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, username string) UserOut
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type UserService struct {
	storage storage.Userer
	logger  *zap.Logger
}

func (u *UserService) GetByLogin(ctx context.Context, username string) UserOut {
	//TODO implement me
	panic("implement me")
}

func NewUserService(storage storage.Userer, logger *zap.Logger) *UserService {
	return &UserService{storage: storage, logger: logger}
}

func (u *UserService) Create(ctx context.Context, user models.User) int {

	_, err := u.storage.Create(ctx, &user)
	if err != nil {
		if v, ok := err.(*pq.Error); ok && v.Code == "23505" {
			return errors.UserServiceUserAlreadyExists

		}
		return errors.UserServiceCreateUserErr
	}

	return errors.NoError

}

func (u *UserService) Update(ctx context.Context, user models.User) error {
	panic("implement me")
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := u.storage.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("user: GetByEmail err", zap.Error(err))
		return nil, err
	}

	return user, nil
}
