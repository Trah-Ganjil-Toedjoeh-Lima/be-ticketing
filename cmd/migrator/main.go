package main

import (
	"fmt"
	"github.com/frchandra/gmcgo/injector"
)

func main() {
	fmt.Println("running migrations")
	migrator := injector.InitializeMigrator()
	migrator.RunMigration("migrate:fresh")
	//TODO: minimizing open too many connection to database by improving query with eager loading
	//TODO: create logger
}
