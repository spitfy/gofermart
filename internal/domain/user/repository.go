package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/repository"
)

var (
	ErrExistsUser   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

type Store struct {
	pool *pgxpool.Pool
}

//go:generate mockgen -destination=storer_mock.go -package=user github.com/spitfy/gofermart/internal/domain/user Storer
type Storer interface {
	RegisterUser(ctx context.Context, user User) (int, error)
	PassByLogin(ctx context.Context, login string) (AuthUser, error)
	Balance(ctx context.Context, userID int) (Balance, error)
	UserBalance(ctx context.Context, userID int) (float64, error)
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		pool: db.Pool(),
	}
}

func (s *Store) RegisterUser(ctx context.Context, user User) (id int, err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		err = errors.Join(err, tx.Commit(ctx))
	}()

	err = tx.QueryRow(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		user.Login, user.Password,
	).Scan(&id)

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		return 0, fmt.Errorf("user already exists: %w", ErrExistsUser)
	case err != nil:
		return 0, fmt.Errorf("failed to insert user: %w", err)
	default:
		return id, nil
	}
}

func (s *Store) PassByLogin(ctx context.Context, login string) (AuthUser, error) {
	var id int
	var password string
	err := s.pool.QueryRow(ctx,
		`SELECT id, password FROM users WHERE login = $1`,
		login,
	).Scan(&id, &password)

	if err != nil {
		return AuthUser{}, fmt.Errorf("failed to get user by login '%s': %w", login, err)
	}
	return AuthUser{
		ID:       id,
		Login:    login,
		Password: password,
	}, nil
}

func (s *Store) UserBalance(ctx context.Context, userID int) (float64, error) {
	var balance float64
	err := s.pool.QueryRow(ctx,
		`SELECT balance FROM users WHERE id = $1`,
		userID,
	).Scan(&balance)
	return balance, err
}

func (s *Store) Balance(ctx context.Context, userID int) (Balance, error) {
	var balance Balance

	err := s.pool.QueryRow(
		ctx,
		`SELECT 
            (SELECT balance FROM users WHERE id = $1), 
            (SELECT COALESCE(SUM(amount), 0) FROM withdrawals WHERE user_id = $1)`,
		userID,
	).Scan(&balance.Current, &balance.Withdrawn)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return Balance{}, fmt.Errorf("user balance not found (userID: %d): %w", userID, ErrUserNotFound)
	case err != nil:
		return Balance{}, fmt.Errorf("failed to get balance for user %d: %w", userID, err)
	default:
		return balance, nil
	}
}
