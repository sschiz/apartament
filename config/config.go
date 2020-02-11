package config

import "github.com/spf13/viper"

func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/apartment")
	viper.AddConfigPath("$HOME/.apartment")

	return viper.ReadInConfig()
}
