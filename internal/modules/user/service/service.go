package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/models"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	"github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"go.uber.org/zap"
)

type Userer interface {
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, username string) (*models.UserDTO, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetUserInfo(ctx context.Context, id int) (*message.UserInfo, error)
}

type UserService struct {
	storage storage.Userer
	logger  *zap.Logger
}

func (u *UserService) GetUserInfo(ctx context.Context, id int) (*message.UserInfo, error) {
	userInfo := message.UserInfo{}
	user, err := u.storage.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("user: GetByID err", zap.Error(err))
		return nil, err
	}

	userInfo.Parser = message.Parser{
		Type:  user.ParserType.String,
		Token: user.ParserToken.String,
	}
	userInfo.Locale = user.Locale.String

	return &userInfo, err

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
	userDB, err := u.storage.GetByID(ctx, user.ID)
	if err != nil {
		u.logger.Error("user: GetUserInfo err", zap.Error(err))
		return err
	}

	userDB.ParserToken.String = user.ParserToken
	userDB.ParserToken.Valid = true
	userDB.ParserType.String = user.ParserType
	userDB.ParserType.Valid = true
	userDB.Locale.String = user.Locale
	userDB.Locale.Valid = true

	err = u.storage.Update(ctx, userDB)
	if err != nil {
		u.logger.Error("user: Update err", zap.Error(err))
		return err
	}

	return nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := u.storage.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("user: GetByEmail err", zap.Error(err))
		return nil, err
	}

	us := &models.User{
		ID:          user.ID,
		Username:    user.Username,
		Password:    user.Hash,
		ParserToken: user.ParserToken.String,
		ParserType:  user.ParserType.String,
		Locale:      user.Locale.String,
	}

	return us, nil
}

//func (u *UserService) GetUserInfo(ctx context.Context, id int) (*models.User, error) {
//	user, err := u.storage.GetUserInfo(ctx, id)
//	if err != nil {
//		u.logger.Error("user: GetByEmail err", zap.Error(err))
//		return nil, err
//	}
//
//	us := &models.User{
//		ID:          user.ID,
//		Username:    user.Username,
//		Password:    user.Hash,
//		ParserToken: user.ParserToken.String,
//		ParserType:  user.ParserType.String,
//		Locale:      user.Locale.String,
//	}
//
//	return us, nil
//}
