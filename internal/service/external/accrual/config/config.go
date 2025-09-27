package config

type Config struct {
	SystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Interval      int    `env:"ACCRUAL_SEND_INTERVAL"`
}
