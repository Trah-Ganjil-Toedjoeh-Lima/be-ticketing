package main

import (
	"fmt"
	"github.com/frchandra/gmcgo/injector"
)

func main() {
	fmt.Println("running migrations")
	migrator := injector.InitializeMigrator()
	migrator.RunMigration("migrate:fresh")
}
