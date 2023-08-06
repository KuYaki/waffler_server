package db

import (
	"context"
	"github.com/KuYaki/waffler_server/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewSqlDB(dbConf config.AppConf) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), dbConf.DatabaseURL)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
