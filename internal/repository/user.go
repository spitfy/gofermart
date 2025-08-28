package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/model"
)

var (
	ErrExistsUser = errors.New("user already exists")
)

func (s *DBStore) RegisterUser(ctx context.Context, user model.User) (int, error) {
	var id int
	err := s.conn.QueryRow(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		user.Login, user.Password,
	).Scan(&id)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		return 0, ErrExistsUser
	case err != nil:
		return 0, err
	default:
		return id, nil
	}
}

func (s *DBStore) PassByLogin(ctx context.Context, login string) (model.AuthUser, error) {
	var id int
	var password string
	err := s.conn.QueryRow(ctx,
		`SELECT id, password FROM users WHERE login = $1`,
		login,
	).Scan(&id, &password)
	if err != nil {
		return model.AuthUser{}, err
	}
	return model.AuthUser{
		ID:       id,
		Login:    login,
		Password: password,
	}, nil
}
