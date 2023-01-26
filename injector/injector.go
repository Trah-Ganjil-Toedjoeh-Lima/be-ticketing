//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/gmcgo/app"
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/frchandra/gmcgo/app/middleware"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/config"
	"github.com/frchandra/gmcgo/database"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	controller.NewUserController,
	middleware.NewUserMiddleware,
)

var ReservationSet = wire.NewSet(
	repository.NewReservationRepository,
	service.NewReservationService,
	controller.NewReservationController,
)

var SeatSet = wire.NewSet(
	repository.NewSeatRepository,
	service.NewSeatService,
)

var TransactionSet = wire.NewSet(
	controller.NewTransactionController,
	repository.NewTransactionRepository,
	service.NewTransactionService,
)

var UtilSet = wire.NewSet(
	util.NewTokenUtil,
)

func InitializeServer() *gin.Engine {
	wire.Build(
		config.NewAppConfig,
		app.NewDatabase,
		app.NewCache,
		UtilSet,
		UserSet,
		SeatSet,
		ReservationSet,
		TransactionSet,
		app.NewRouter,
	)
	return nil
}

func InitializeMigrator() *database.Migrator {
	wire.Build(
		config.NewAppConfig,
		app.NewDatabase,
		database.NewMigrator,
	)
	return nil
}
