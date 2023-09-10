package storages

import (
	ustorage "github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"gorm.io/gorm"
)

type Storages struct {
	User    ustorage.Userer
	Waffler storage.WafflerStorager
}

func NewStorages(conn *gorm.DB) *Storages {
	return &Storages{
		User:    ustorage.NewUserStorage(conn),
		Waffler: storage.NewWafflerStorage(conn),
	}
}
