package service

import (
	"context"
	"github.com/KuYaki/waffler_server/internal/modules/message"
)

type Auther interface {
	Register(ctx context.Context, username, password string) (int, error)
	Login(ctx context.Context, user message.User) (*AuthorizeOut, int, error)
	AuthorizeRefresh(ctx context.Context, idUser int) (*AuthorizeOut, error)
}

type AuthorizeOut struct {
	UserID       int
	AccessToken  string
	RefreshToken string
}
