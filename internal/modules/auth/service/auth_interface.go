package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/models"
)

type Auther interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, user models.User) (*AuthorizeOut, error)
	AuthorizeRefresh(ctx context.Context, idUser int) (*AuthorizeOut, error)
}

type AuthorizeOut struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}
