package storages

import (
	bstorage "github.com/KuYaki/waffler_server/internal/modules/bot_translator/storage"
	ustorage "github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"gorm.io/gorm"
)

type Storages struct {
	User       ustorage.Userer
	Waffler    storage.WafflerStorager
	BotStorage bstorage.BotStorager
}

func NewStorages(conn *gorm.DB) *Storages {
	return &Storages{
		User:       ustorage.NewUserStorage(conn),
		Waffler:    storage.NewWafflerStorage(conn),
		BotStorage: bstorage.NewWafflerStorage(conn),
	}
}
