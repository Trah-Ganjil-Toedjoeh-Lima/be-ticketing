// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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

// Injectors from injector.go:

func InitializeServer() *gin.Engine {
	appConfig := config.NewAppConfig()
	client := app.NewCache(appConfig)
	tokenUtil := util.NewTokenUtil(client, appConfig)
	logger := app.NewLogger(appConfig)
	userMiddleware := middleware.NewUserMiddleware(tokenUtil, logger)
	db := app.NewDatabase(appConfig, logger)
	logUtil := util.NewLogUtil(logger)
	userRepository := repository.NewUserRepository(db, logUtil)
	userService := service.NewUserService(userRepository, tokenUtil)
	adminMiddleware := middleware.NewAdminMiddleware(tokenUtil, logger, appConfig, userService)
	gateMiddleware := middleware.NewGateMiddleware(appConfig)
	userController := controller.NewUserController(userService, tokenUtil, appConfig)
	transactionRepository := repository.NewTransactionRepository(db, logUtil)
	seatRepository := repository.NewSeatRepository(db, logUtil)
	transactionService := service.NewTransactionService(transactionRepository, userRepository, seatRepository, appConfig)
	reservationService := service.NewReservationService(appConfig, transactionService)
	seatService := service.NewSeatService(appConfig, seatRepository, transactionRepository)
	reservationController := controller.NewReservationController(appConfig, db, logUtil, reservationService, transactionService, seatService)
	snapUtil := util.NewSnapUtil(appConfig)
	transactionController := controller.NewTransactionController(transactionService, userService, snapUtil, logUtil)
	emailUtil := util.NewEmailUtil(appConfig)
	eTicketUtil := util.NewETicketUtil(appConfig)
	snapService := service.NewSnapService(transactionService, seatService, transactionRepository, snapUtil, emailUtil, eTicketUtil)
	snapController := controller.NewSnapController(snapService, snapUtil, transactionService, logUtil)
	gateController := controller.NewGateController(appConfig)
	engine := app.NewRouter(userMiddleware, adminMiddleware, gateMiddleware, userController, reservationController, transactionController, snapController, gateController)
	return engine
}

func InitializeMigrator() *database.Migrator {
	appConfig := config.NewAppConfig()
	logger := app.NewLogger(appConfig)
	db := app.NewDatabase(appConfig, logger)
	migrator := database.NewMigrator(db)
	return migrator
}

func InitializeEmail() *util.EmailUtil {
	appConfig := config.NewAppConfig()
	emailUtil := util.NewEmailUtil(appConfig)
	return emailUtil
}

// injector.go:

var MiddlewareSet = wire.NewSet(middleware.NewUserMiddleware, middleware.NewAdminMiddleware, middleware.NewGateMiddleware)

var UserSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, controller.NewUserController)

var ReservationSet = wire.NewSet(service.NewReservationService, controller.NewReservationController)

var SeatSet = wire.NewSet(repository.NewSeatRepository, service.NewSeatService)

var TransactionSet = wire.NewSet(controller.NewTransactionController, repository.NewTransactionRepository, service.NewTransactionService)

var SnapSet = wire.NewSet(controller.NewSnapController, service.NewSnapService)

var GateSet = wire.NewSet(controller.NewGateController)

var UtilSet = wire.NewSet(util.NewTokenUtil, util.NewSnapUtil, util.NewEmailUtil, util.NewETicketUtil, util.NewLogUtil)
