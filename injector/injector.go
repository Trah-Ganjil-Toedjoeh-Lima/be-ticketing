//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/frchandra/ticketing-gmcgo/app"
	"github.com/frchandra/ticketing-gmcgo/app/controller"
	"github.com/frchandra/ticketing-gmcgo/app/middleware"
	"github.com/frchandra/ticketing-gmcgo/app/repository"
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/frchandra/ticketing-gmcgo/database"
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

var SnapSet = wire.NewSet(
	controller.NewSnapController,
	service.NewSnapService,
)

var UtilSet = wire.NewSet(
	util.NewTokenUtil,
	util.NewSnapUtil,
	util.NewEmailUtil,
	util.NewETicketUtil,
)

func InitializeServer() *gin.Engine {
	wire.Build(
		config.NewAppConfig,
		app.NewDatabase,
		app.NewCache,
		app.NewLogger,
		UtilSet,
		UserSet,
		SeatSet,
		ReservationSet,
		TransactionSet,
		SnapSet,
		app.NewRouter,
	)
	return nil
}

func InitializeMigrator() *database.Migrator {
	wire.Build(
		config.NewAppConfig,
		app.NewLogger,
		app.NewDatabase,
		database.NewMigrator,
	)
	return nil
}

func InitializeEmail() *util.EmailUtil {
	wire.Build(
		config.NewAppConfig,
		util.NewEmailUtil,
	)
	return nil
}
