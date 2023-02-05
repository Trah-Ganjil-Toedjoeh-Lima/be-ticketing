package main

import (
	"fmt"
	"github.com/frchandra/ticketing-gmcgo/injector"
)

func main() {
	fmt.Println("running migrations")
	migrator := injector.InitializeMigrator()
	migrator.RunMigration("migrate:fresh")
}
