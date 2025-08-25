package service

import "github.com/spitfy/gofermart/internal/config"

type UserService struct {
	cfg *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		cfg: cfg,
	}
}

func (us *UserService) RegisterUser() error {
	return nil
}

func (us *UserService) LoginUser() error {
	return nil
}
