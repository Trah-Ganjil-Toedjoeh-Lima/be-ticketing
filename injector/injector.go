//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/database"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

/*var userSet = wire.NewSet(
	controller.NewUserController,
	service.NewUserService,
	repository.NewUserRepository,
)

func InitializeServer() *app.Server {
	wire.Build(
		config.NewAppConfig,
		app.NewServer,
	)
	return nil
}*/

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
		config.NewAppConfig,
		database.NewMigration,
		database.NewMigrator,
	)
	return nil
}

/*func InitializeUserController(db *gorm.DB) *controller.UserController {
	wire.Build(
		controller.NewUserController,
		service.NewUserService,
		repository.NewUserRepository,
	)
	return nil
}*/
