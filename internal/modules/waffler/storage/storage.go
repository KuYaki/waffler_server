package storage

import "github.com/jackc/pgx/v5/pgxpool"

type WafflerStorager interface {
	Select()
}

type WafflerStorage struct {
	conn *pgxpool.Pool
}

func NewWafflerStorage(conn *pgxpool.Pool) WafflerStorager {
	return &WafflerStorage{conn: conn}
}
func (s WafflerStorage) Select() {

}
