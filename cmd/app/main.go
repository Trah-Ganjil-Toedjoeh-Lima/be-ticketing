package main

import (
	injector "github.com/frchandra/gmcgo/injector"
)

func main() {
	server := injector.InitializeServer()
	server.Run()
}
