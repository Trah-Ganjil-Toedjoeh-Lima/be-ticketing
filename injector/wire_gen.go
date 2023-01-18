// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/database"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from injector.go:

func InitializeServer() *gin.Engine {
	db := app.NewDatabase()
	userRepository := repository.NewUserRepository(db)
	client := app.NewCache()
	tokenUtil := util.NewTokenUtil(client)
	userService := service.NewUserService(userRepository, tokenUtil)
	userController := controller.NewUserController(userService, tokenUtil)
	engine := app.NewRouter(userController)
	return engine
}

func InitializeMigrator() *database.Migrator {
	db := app.NewDatabase()
	migration := database.NewMigration()
	migrator := database.NewMigrator(db, migration)
	return migrator
}

// injector.go:

var UserSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, controller.NewUserController)

var UtilSet = wire.NewSet(util.NewTokenUtil)
