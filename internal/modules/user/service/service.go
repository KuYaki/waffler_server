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
	Create(ctx context.Context, user *models.UserDTO) error
	Update(ctx context.Context, user message.UserInfo, idUser int) error
	GetByLogin(ctx context.Context, username string) (*models.UserDTO, error)
	GetByID(ctx context.Context, id int) (*message.User, error)
	GetUserInfo(ctx context.Context, id int) (*message.UserInfo, error)
	ExistsUser(ctx context.Context, username string) (bool, error)
}

type UserService struct {
	storage storage.Userer
	logger  *zap.Logger
}

func (u *UserService) GetUserInfo(ctx context.Context, id int) (*message.UserInfo, error) {
	_, err := u.storage.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("user: GetByID err", zap.Error(err))
		return nil, err
	}

	//  maintain compatibility  ToDo: change it
	userInfo := &message.UserInfo{
		Parser: message.Parser{
			Type:  models.NotFound,
			Token: "not found",
		},
		Locale: "not found",
	}

	return userInfo, err

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

func (u *UserService) Create(ctx context.Context, user *models.UserDTO) error {
	err := u.storage.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil

}

func (u *UserService) Update(ctx context.Context, user message.UserInfo, idUser int) error {

	userDB, err := newUserDTO(user)
	if err != nil {
		u.logger.Error("user: newUserDTO err", zap.Error(err))
		return err
	}
	userDB.ID = idUser

	err = u.storage.Update(ctx, userDB)
	if err != nil {
		u.logger.Error("user: Update err", zap.Error(err))
		return err
	}

	return nil
}

func (u *UserService) GetByID(ctx context.Context, id int) (*message.User, error) {
	user, err := u.storage.GetByID(ctx, id)
	if err != nil {
		u.logger.Error("user: GetByID err", zap.Error(err))
		return nil, err
	}

	us := &message.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.PwdHash,
	}

	return us, nil
}

func (u *UserService) ExistsUser(ctx context.Context, username string) (bool, error) {
	user, err := u.storage.UserExists(ctx, username)
	if err != nil {
		u.logger.Error("user: GetByUsername err", zap.Error(err))
		return user, err
	}

	return user, err
}

func newUserDTO(user message.UserInfo) (*models.UserDTO, error) {
	//  maintain compatibility  ToDo: change it
	userDB := models.UserDTO{}
	//if user.Parser.Token != "" {
	//	err = userDB.ParserToken.Scan(user.Parser.Token)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//if user.Parser.Type != "" {
	//	err = userDB.ParserType.Scan(user.Parser.Type)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//if user.Locale != "" {
	//	err = userDB.Locale.Scan(user.Locale)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return &userDB, nil
}
