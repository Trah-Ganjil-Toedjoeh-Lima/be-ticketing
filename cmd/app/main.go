package main

import (
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/injector"
)

func main() {
	appConfig := config.NewAppConfig()
	router := injector.InitializeServer()
	router.Run(":" + appConfig.AppPort)

	//TODO: simulasi runtime error, apakah seluruh aplikasi berhenti?
	//TODO: edge case => menghijaukan kursi kuning yang di cim oleh user yg tidak bertanggunjawab
}
