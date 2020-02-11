package main

import (
	"github.com/spf13/viper"
	"github.com/sschiz/apartament/config"
	"github.com/sschiz/apartament/internal/server"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("server.port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
