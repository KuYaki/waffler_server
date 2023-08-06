package service

import "github.com/KuYaki/waffler_server/internal/models"

type UserOut struct {
	User      *models.User `json:"user"`
	ErrorCode int          `json:"error_code"`
}
