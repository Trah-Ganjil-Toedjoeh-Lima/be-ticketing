package main

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/frchandra/ticketing-gmcgo/injector"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.NewAppConfig()

	if appConfig.IsProduction == true {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := injector.InitializeServer()
	router.Run(":" + appConfig.AppPort)

	//TODO: simulasi runtime error, apakah seluruh aplikasi berhenti?
	//TODO: cari cara lain untuk menormalkan transaksi ngambang dari user tak bertanggungjawab menggunakan queue spt redis, kafka, dll? (sekarang pakai timestamp)
	//TODO: minimizing open too many connection to database by improving query with eager loading
	//TODO: create delete seat transaction feature
	//TODO: index database column to increase speed, learn how db index works and pros/cons (not necessary yet)
	//TODO: add minio service
	//TODO: add migration entrypoint

}
