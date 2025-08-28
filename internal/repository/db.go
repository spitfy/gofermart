package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/migration"
)

type DBStore struct {
	conf *config.Config
	conn *pgxpool.Pool
}

func NewDBStore(conf *config.Config) (*DBStore, error) {
	if err := migrate(conf); err != nil {
		return nil, err
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conf.DB.DatabaseURI)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &DBStore{
		conf,
		pool,
	}, nil
}

func migrate(conf *config.Config) error {
	m := migration.NewMigration(conf)
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Close() {
	s.conn.Close()
}
