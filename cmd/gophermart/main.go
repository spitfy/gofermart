package main

import (
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/handler"
	"github.com/spitfy/gofermart/internal/repository"
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

	store, err := repository.NewDBStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()
	us := service.NewUserService(cfg, store)

	return handler.Serve(cfg, us)
}
