package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/spitfy/gofermart/internal/config/db"
	handler "github.com/spitfy/gofermart/internal/handler/config"
	"github.com/spitfy/gofermart/internal/helper"
	auth "github.com/spitfy/gofermart/internal/middleware/auth/config"
	loggerConf "github.com/spitfy/gofermart/internal/middleware/logger/config"
	accrual "github.com/spitfy/gofermart/internal/service/external/accrual/config"
)

type Config struct {
	DB      db.Config
	Handler handler.Config
	Auth    auth.Config
	Accrual accrual.Config
	Logger  loggerConf.Config
}

var (
	SecretKey       = "SecRetKey"
	AccrualInterval = 1
	maxTries        = uint(3)
)

type Default struct {
	RunAddress      string
	AccrualAddress  string
	AccrualInterval int
	DatabaseURI     string
	LogLevel        string
	Secret          string
	MaxTries        uint
}

func newDefault() (*Default, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	moduleRoot, err := helper.FindModuleRoot(wd)
	if err != nil {
		return nil, err
	}
	envPath := filepath.Join(moduleRoot, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	viper.AutomaticEnv()

	d := Default{
		RunAddress:      viper.GetString("run_address"),
		AccrualAddress:  viper.GetString("accrual_system_address"),
		AccrualInterval: viper.GetInt("accrual_send_interval"),
		DatabaseURI:     viper.GetString("database_uri"),
		LogLevel:        viper.GetString("log_level"),
		Secret:          viper.GetString("secret"),
		MaxTries:        viper.GetUint("max_tries"),
	}
	if d.Secret == "" {
		d.Secret = SecretKey
	}
	if d.AccrualInterval == 0 {
		d.AccrualInterval = AccrualInterval
	}
	if d.MaxTries == 0 {
		d.MaxTries = maxTries
	}
	return &d, nil
}

func GetConfig() *Config {
	d, err := newDefault()
	if err != nil {
		log.Println(err)
	}

	conf := &Config{
		Auth: auth.Config{
			SecretKey: d.Secret,
		},
		Accrual: accrual.Config{
			Interval: d.AccrualInterval,
		},
		DB: db.Config{
			MaxRetries: d.MaxTries,
		},
	}

	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}
	if conf.Handler.RunAddress == "" {
		conf.Handler.RunAddress = d.RunAddress
	}
	if conf.DB.DatabaseURI == "" {
		conf.DB.DatabaseURI = d.DatabaseURI
	}
	if conf.Accrual.SystemAddress == "" {
		conf.Accrual.SystemAddress = d.AccrualAddress
	}
	flag.StringVar(&conf.Handler.RunAddress, "a", conf.Handler.RunAddress, "server address")
	flag.StringVar(&conf.DB.DatabaseURI, "d", conf.DB.DatabaseURI, "database DSN address")
	flag.StringVar(&conf.Accrual.SystemAddress, "r", conf.Accrual.SystemAddress, "accrual server address")
	flag.StringVar(&conf.Logger.LogLevel, "l", d.LogLevel, "Logger level")

	flag.Parse()

	return conf
}
