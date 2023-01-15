package main

import (
	"github.com/frchandra/gmcgo/app"
	injector "github.com/frchandra/gmcgo/injector"
)

func main() {
	var server *app.Server = injector.InitializeServer()
	server.Run()
}
