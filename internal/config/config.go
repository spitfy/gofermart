package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/spitfy/gofermart/internal/config/db"
	handler "github.com/spitfy/gofermart/internal/handler/config"
	"log"
)

type Config struct {
	DB      db.Config
	Handler handler.Config
}

const (
	DefaultRunAddress string = ":8080"
)

func GetConfig() *Config {
	conf := &Config{}

	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}
	if conf.Handler.RunAddress == "" {
		conf.Handler.RunAddress = DefaultRunAddress
	}
	flag.StringVar(&conf.Handler.RunAddress, "a", conf.Handler.RunAddress, "server address")
	flag.StringVar(&conf.DB.DatabaseURI, "d", conf.DB.DatabaseURI, "database DSN address")

	flag.Parse()

	return conf
}
