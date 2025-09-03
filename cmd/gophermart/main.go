package main

import (
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/middleware/logger"
	"github.com/spitfy/gofermart/internal/repository"
	"github.com/spitfy/gofermart/internal/repository/order"
	"github.com/spitfy/gofermart/internal/repository/user"
	serviceAccrual "github.com/spitfy/gofermart/internal/service/external/accrual"
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

	//balanceStore := balance.NewStore(store)
	orderStore := order.NewStore(store)
	as := serviceAccrual.NewService(cfg, orderStore)
	os := serviceOrder.NewService(cfg, orderStore, as)

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}

	s := handler.Service{
		Auth:         auth.New(cfg.Auth.SecretKey),
		UserService:  us,
		OrderService: os,
		Logger:       l,
	}

	return handler.Serve(cfg, s)
}
