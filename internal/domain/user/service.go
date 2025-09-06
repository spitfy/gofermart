package user

import (
	"context"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cfg *config.Config
	s   Storer
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (us *Service) RegisterUser(ctx context.Context, user User) (int, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash
	return us.s.RegisterUser(ctx, user)
}

func (us *Service) LoginUser(ctx context.Context, user User) (int, error) {
	u, err := us.s.PassByLogin(ctx, user.Login)
	if err != nil {
		return -1, err
	}
	if !checkPasswordHash(user.Password, u.Password) {
		return -1, auth.ErrUnAuth
	}
	return u.ID, nil
}

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (us *Service) Balance(ctx context.Context, userID int) (model.Balance, error) {
	return us.s.Balance(ctx, userID)
}

func (us *Service) CanWithdraw(ctx context.Context, userID int, w float64) (bool, error) {
	b, err := us.s.UserBalance(ctx, userID)
	if err != nil {
		return false, err
	}
	return b >= w, nil
}
