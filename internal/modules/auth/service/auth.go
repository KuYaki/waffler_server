package service

import (
	"context"
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/component"
	"github.com/KuYaki/waffler_server/internal/infrastructure/errors"
	"github.com/KuYaki/waffler_server/internal/infrastructure/tools/cryptography"
	"github.com/KuYaki/waffler_server/internal/models"
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

func (a *Auth) Login(ctx context.Context, user models.User) *AuthorizeOut {
	// 1. получаем юзера по username
	userDb := a.user.GetByLogin(ctx, user.Username)
	if userDb.ErrorCode != errors.NoError {
		return &AuthorizeOut{
			ErrorCode: userDb.ErrorCode,
		}
	}
	// 2. проверяем пароль
	if !cryptography.CheckPassword(user.Hash, userDb.User.Password) {
		return &AuthorizeOut{
			ErrorCode: errors.AuthServiceWrongPasswordErr,
		}
	}
	user.ID = userDb.User.ID

	// 3. генерируем токены
	accessToken, refreshToken, err := a.generateTokens(&user)
	if err != nil {
		return &AuthorizeOut{
			ErrorCode: errors.AuthServiceTokenGenerationErr,
		}
	}
	// 4. возвращаем токены
	return &AuthorizeOut{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *Auth) Register(ctx context.Context, username, password string) (int, int) {
	hashPass, err := cryptography.HashPassword(password)
	if err != nil {
		return http.StatusInternalServerError, errors.HashPasswordError
	}
	dto := models.User{
		Username: username,
		Hash:     hashPass,
	}

	userOut := a.user.Create(ctx, dto)
	if userOut != errors.NoError {
		if userOut == errors.UserServiceUserAlreadyExists {
			return http.StatusConflict, userOut

		}
		return http.StatusInternalServerError, userOut
	}

	return http.StatusOK, errors.NoError
}

func (a *Auth) AuthorizeRefresh(ctx context.Context, idUser int) *AuthorizeOut {
	userOut, err := a.user.GetByID(ctx, idUser)
	if err != nil {
		return &AuthorizeOut{
			ErrorCode: errors.AuthServiceGeneralErr,
		}
	}

	accessToken, refreshToken, err := a.generateTokens(userOut)
	if err != nil {
		return &AuthorizeOut{
			ErrorCode: errors.AuthServiceTokenGenerationErr,
		}
	}

	return &AuthorizeOut{
		UserID:       idUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *Auth) generateTokens(user *models.User) (string, string, error) {
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
