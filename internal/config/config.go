package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/spitfy/gofermart/internal/config/db"
	handler "github.com/spitfy/gofermart/internal/handler/config"
	auth "github.com/spitfy/gofermart/internal/middleware/auth/config"
	loggerConf "github.com/spitfy/gofermart/internal/middleware/logger/config"
	accrual "github.com/spitfy/gofermart/internal/service/external/accrual/config"
	"log"
)

type Config struct {
	DB      db.Config
	Handler handler.Config
	Auth    auth.Config
	Accrual accrual.Config
	Logger  loggerConf.Config
}

const (
	DefaultRunAddress     string = ":8080"
	DefaultAccrualAddress string = "http://localhost:8082"
	DefaultDatabaseURI    string = "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable"
	SecretKey             string = "**SecRetKey#!45**"
	DefaultLogLevel       string = "info"
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
	flag.StringVar(&conf.Logger.LogLevel, "l", DefaultLogLevel, "Logger level")

	flag.Parse()

	return conf
}
