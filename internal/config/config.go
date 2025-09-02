package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	auth "github.com/spitfy/gofermart/internal/auth/config"
	"github.com/spitfy/gofermart/internal/config/db"
	handler "github.com/spitfy/gofermart/internal/handler/config"
	accrual "github.com/spitfy/gofermart/internal/service/external/accrual/config"
	"log"
)

type Config struct {
	DB      db.Config
	Handler handler.Config
	Auth    auth.Config
	Accrual accrual.Config
}

const (
	DefaultRunAddress     string = ":8080"
	DefaultAccrualAddress string = "http://localhost:8082"
	DefaultDatabaseURI    string = "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable"
	SecretKey             string = "**SecRetKey#!45**"
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
	if conf.Accrual.SystemAddress == "" {
		conf.Accrual.SystemAddress = DefaultAccrualAddress
	}
	flag.StringVar(&conf.Handler.RunAddress, "a", conf.Handler.RunAddress, "server address")
	flag.StringVar(&conf.DB.DatabaseURI, "d", conf.DB.DatabaseURI, "database DSN address")
	flag.StringVar(&conf.Accrual.SystemAddress, "r", conf.Accrual.SystemAddress, "accrual server address")

	flag.Parse()

	return conf
}
