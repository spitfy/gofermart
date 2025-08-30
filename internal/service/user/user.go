package user

import (
	"context"
	"github.com/spitfy/gofermart/internal/auth"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	cfg *config.Config
	s   Storer
}

type Storer interface {
	RegisterUser(ctx context.Context, user model.User) (int, error)
	PassByLogin(ctx context.Context, login string) (model.AuthUser, error)
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (us *Service) RegisterUser(ctx context.Context, user model.User) (int, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash
	return us.s.RegisterUser(ctx, user)
}

func (us *Service) LoginUser(ctx context.Context, user model.User) (int, error) {
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
