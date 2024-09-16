package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
)

func NewPostgres() (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(viper.GetString("DB_URI"))
	if err != nil {
		return nil, err
	}

	// UUID support for google/uuid
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
