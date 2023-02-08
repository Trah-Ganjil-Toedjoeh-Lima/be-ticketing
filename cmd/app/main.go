package main

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/frchandra/ticketing-gmcgo/injector"
)

func main() {
	appConfig := config.NewAppConfig()
	router := injector.InitializeServer()
	router.Run(":" + appConfig.AppPort)

	//TODO: simulasi runtime error, apakah seluruh aplikasi berhenti?
	//TODO: cari cara lain untuk menormalkan transaksi ngambang dari user tak bertanggungjawab menggunakan queue spt redis, kafka, dll? (sekarang pakai timestamp)
	//TODO: minimizing open too many connection to database by improving query with eager loading
	//TODO: create delete seat transaction feature
	//TODO: move picture ticket to docker container
	//TODO: pretify transaction details data for FE
}
