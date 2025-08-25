package main

import (
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/service"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	cfg := config.GetConfig()
	userService := service.NewUserService(cfg)

	/*store, err := repository.CreateStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()
	s := service.NewService(*cfg, store)

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}*/

	return handler.Serve(cfg, userService)
}
