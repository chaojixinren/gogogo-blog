package main

import (
	"gogogo/config"
	"gogogo/router"
)

func main() {
	config.InitConfig()
	server := router.SetupRouter()

	server.Run(config.AppConfig.App.Port)
}
