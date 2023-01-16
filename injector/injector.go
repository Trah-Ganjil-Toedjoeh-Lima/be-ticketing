//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/database"
	"github.com/google/wire"
)

func InitializeServer() *app.Server {
	wire.Build(
		config.NewAppConfig,
		app.NewServer,
	)
	return nil
}

func InitializeMigrator() *database.Migrator {
	wire.Build(
		config.NewAppConfig,
		database.NewMigration,
		database.NewMigrator,
	)
	return nil
}
