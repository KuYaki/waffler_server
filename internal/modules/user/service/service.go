package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"go.uber.org/zap"
)

type Userer interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, username string) (*models.UserDTO, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type UserService struct {
	storage storage.Userer
	logger  *zap.Logger
}

func (u *UserService) GetByLogin(ctx context.Context, username string) (*models.UserDTO, error) {
	user, err := u.storage.GetByUsername(ctx, username)
	if err != nil {
		u.logger.Error("user: GetByUsername err", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func NewUserService(storage storage.Userer, components *component.Components) *UserService {
	return &UserService{storage: storage, logger: components.Logger}
}

func (u *UserService) Create(ctx context.Context, user models.User) error {
	us := &models.UserDTO{
		Username: user.Username,
		Hash:     user.Password,
	}
	err := u.storage.Create(ctx, us)
	if err != nil {
		return err
	}

	return nil

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

	us := &models.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Hash,
		TokenGPT: user.TokenGPT,
	}

	return us, nil
}
