//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/database"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	controller.NewUserController,
)

func InitializeRouter() *gin.Engine {
	wire.Build(
		app.NewDatabase,
		UserSet,
		app.NewRouter,
	)
	return nil
}

func InitializeMigrator() *database.Migrator {
	wire.Build(
		app.NewDatabase,
		database.NewMigration,
		database.NewMigrator,
	)
	return nil
}
