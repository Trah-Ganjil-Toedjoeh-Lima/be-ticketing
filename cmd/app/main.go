package main

import (
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/injector"
)

func main() {
	appConfig := config.NewAppConfig()
	router := injector.InitializeRouter()
	router.Run(":" + appConfig.AppPort)
}
