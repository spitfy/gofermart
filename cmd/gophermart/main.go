package main

import (
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/domain/accrual"
	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/middleware/logger"
	"github.com/spitfy/gofermart/internal/repository"
	serviceAccrual "github.com/spitfy/gofermart/internal/service/external/accrual"
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
	us := user.NewService(cfg, userStore)

	accrualStore := accrual.NewStore(store)
	as := accrual.NewService(cfg, accrualStore)

	orderStore := order.NewStore(store)
	extAs := serviceAccrual.NewService(cfg, as)
	os := order.NewService(cfg, orderStore, extAs)

	withdrawStore := withdraw.NewStore(store)
	ws := withdraw.NewService(cfg, withdrawStore)

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}

	s := handler.Service{
		Auth:            auth.New(cfg.Auth.SecretKey),
		UserService:     us,
		OrderService:    os,
		WithdrawService: ws,
		Logger:          l,
	}

	return handler.Serve(cfg, s)
}
