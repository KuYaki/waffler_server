package storages

import (
	ustorage "github.com/KuYaki/waffler_server/internal/modules/user/storage"
	"github.com/KuYaki/waffler_server/internal/modules/waffler/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storages struct {
	User    ustorage.Userer
	Waffler storage.WafflerStorager
}

func NewStorages(conn *pgxpool.Pool) *Storages {
	return &Storages{
		User:    ustorage.NewUserStorage(conn),
		Waffler: storage.NewWafflerStorage(conn),
	}
}
