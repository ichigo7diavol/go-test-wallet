package config

import (
	"github.com/spf13/viper"
)

const (
	EnvPrefix   = "WALLET_APP"
	DefaultPort = "8080"
)

type Config struct {
	Port string
	Dsn  string
}

func Load() *Config {
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	viper.SetDefault("port", DefaultPort)

	viper.BindEnv("port", "PORT")
	viper.BindEnv("dsn", "DSN")

	port := viper.GetString("port")
	dsn := viper.GetString("dsn")

	return &Config{
		Port: port,
		Dsn:  dsn,
	}
}
