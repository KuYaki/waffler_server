package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/models"
)

type Auther interface {
	Register(ctx context.Context, username, password string) (int, int)
	Login(ctx context.Context, user models.User) *AuthorizeOut
	AuthorizeRefresh(ctx context.Context, idUser int) *AuthorizeOut
}

type AuthorizeOut struct {
	UserID       int
	AccessToken  string
	RefreshToken string
	ErrorCode    int
}
