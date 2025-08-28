package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	auth "github.com/spitfy/gofermart/internal/auth/config"
	"github.com/spitfy/gofermart/internal/config/db"
	handler "github.com/spitfy/gofermart/internal/handler/config"
	"log"
)

type Config struct {
	DB      db.Config
	Handler handler.Config
	Auth    auth.Config
}

const (
	DefaultRunAddress  string = ":8080"
	DefaultDatabaseURI string = "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable"
	SecretKey          string = "**SecRetKey#!45**"
)

func GetConfig() *Config {
	conf := &Config{
		Auth: auth.Config{
			SecretKey: SecretKey,
		},
	}

	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}
	if conf.Handler.RunAddress == "" {
		conf.Handler.RunAddress = DefaultRunAddress
	}
	if conf.DB.DatabaseURI == "" {
		conf.DB.DatabaseURI = DefaultDatabaseURI
	}
	flag.StringVar(&conf.Handler.RunAddress, "a", conf.Handler.RunAddress, "server address")
	flag.StringVar(&conf.DB.DatabaseURI, "d", conf.DB.DatabaseURI, "database DSN address")

	flag.Parse()

	return conf
}
