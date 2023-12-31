package service

import (
	"context"
	"errors"
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/modules/message"
	uservice "github.com/KuYaki/waffler_server/internal/modules/user/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Auth struct {
	conf config.AppConf
	user uservice.Userer

	tokenManager cryptography.TokenManager
	hash         cryptography.Hasher
	logger       *zap.Logger
}

func NewAuthService(user uservice.Userer, components *component.Components) *Auth {
	return &Auth{conf: components.Conf,
		user:         user,
		tokenManager: components.TokenManager,
		hash:         components.Hash,
		logger:       components.Logger,
	}
}

func (a *Auth) Login(ctx context.Context, user message.User) (*AuthorizeOut, int, error) {
	// 1. получаем юзера по username
	userDb, err := a.user.GetByLogin(ctx, user.Username)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}
	// 2. проверяем пароль
	if !cryptography.CheckPassword(userDb.Hash, user.Password) {
		a.logger.Error("user: CheckPassword err", zap.Error(err))
		return nil, http.StatusUnauthorized, errors.New("wrong password")
	}
	user.ID = userDb.ID

	// 3. генерируем токены
	accessToken, refreshToken, err := a.generateTokens(&user)
	if err != nil {
		a.logger.Error("user: generateTokens err", zap.Error(err))
		return nil, http.StatusBadRequest, err
	}
	// 4. возвращаем токены
	return &AuthorizeOut{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, http.StatusOK, nil
}

func (a *Auth) Register(ctx context.Context, username, password string) (int, error) {
	existUser, err := a.user.ExistsUser(ctx, username)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if existUser {
		return http.StatusConflict, errors.New("user already exists")
	}

	hashPass, err := cryptography.HashPassword(password)
	if err != nil {
		return http.StatusBadRequest, err
	}
	dto := message.User{
		Username: username,
		Password: hashPass,
	}

	err = a.user.Create(ctx, dto)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (a *Auth) AuthorizeRefresh(ctx context.Context, idUser int) (*AuthorizeOut, error) {
	userOut, err := a.user.GetByID(ctx, idUser)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := a.generateTokens(userOut)
	if err != nil {
		return nil, err
	}

	return &AuthorizeOut{
		UserID:       idUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *Auth) generateTokens(user *message.User) (string, string, error) {
	accessToken, err := a.tokenManager.CreateToken(
		strconv.Itoa(user.ID),
		a.conf.Token.AccessTTL,
		cryptography.AccessToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", err
	}
	refreshToken, err := a.tokenManager.CreateToken(
		strconv.Itoa(user.ID),
		a.conf.Token.RefreshTTL,
		cryptography.RefreshToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
