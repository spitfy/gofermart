package main

import (
	"github.com/spitfy/gofermart/internal/auth"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/repository"
	"github.com/spitfy/gofermart/internal/repository/order"
	"github.com/spitfy/gofermart/internal/repository/user"
	serviceOrder "github.com/spitfy/gofermart/internal/service/order"
	serviceUser "github.com/spitfy/gofermart/internal/service/user"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	cfg := config.GetConfig()

	store, err := repository.NewDBStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()

	userStore := user.NewStore(store)
	us := serviceUser.NewService(cfg, userStore)
	orderStore := order.NewStore(store)
	os := serviceOrder.NewService(cfg, orderStore)

	s := handler.Service{
		Auth:         auth.New(cfg.Auth.SecretKey),
		UserService:  us,
		OrderService: os,
	}

	return handler.Serve(cfg, s)
}
