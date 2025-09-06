package user

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository"
)

var (
	ErrExistsUser = errors.New("user already exists")
)

type Store struct {
	*repository.DBStore
}

type Storer interface {
	RegisterUser(ctx context.Context, user model.User) (int, error)
	PassByLogin(ctx context.Context, login string) (model.AuthUser, error)
	Balance(ctx context.Context, userID int) (model.Balance, error)
	UserBalance(ctx context.Context, userID int) (float64, error)
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		DBStore: db,
	}
}

func (s *Store) RegisterUser(ctx context.Context, user model.User) (int, error) {
	var id int
	err := s.Conn.QueryRow(ctx,
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

func (s *Store) PassByLogin(ctx context.Context, login string) (model.AuthUser, error) {
	var id int
	var password string
	err := s.Conn.QueryRow(ctx,
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

func (s *Store) UserBalance(ctx context.Context, userID int) (float64, error) {
	var balance float64
	err := s.Conn.QueryRow(ctx,
		`SELECT balance FROM users WHERE id = $1`,
		userID,
	).Scan(&balance)
	return balance, err
}

func (s *Store) Balance(ctx context.Context, userID int) (model.Balance, error) {
	var balance model.Balance
	err := s.Conn.QueryRow(
		ctx,
		`select (SELECT balance FROM users WHERE id = $1), 
       				(select coalesce(SUM(amount), 0) from withdrawals where user_id = $1)`,
		userID,
	).Scan(&balance.Current, &balance.Withdrawn)
	return balance, err
}
