package app

import (
	"github.com/frchandra/ticketing-gmcgo/app/controller"
	"github.com/frchandra/ticketing-gmcgo/app/middleware"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	config *config.AppConfig,

	userMiddleware *middleware.UserMiddleware,
	adminMiddleware *middleware.AdminMiddleware,
	gateMiddleware *middleware.GateMiddleware,
	qrMiddleware *middleware.ScanQrMiddleware,

	homeController *controller.HomeController,
	userController *controller.UserController,
	authController *controller.AuthController,
	reservationController *controller.ReservationController,
	txController *controller.TransactionController,
	snapController *controller.SnapController,
	gateController *controller.ConfigController,
	seatController *controller.SeatController,
) *gin.Engine {
	var router *gin.Engine
	if config.IsProduction == true {
		router = gin.New()
	} else {
		router = gin.Default()
	}

	//Health Check (scope: public)
	router.GET(config.EndpointPrefix+"health", homeController.HealthCheck)

	//Midtrans Webhook (scope: public)
	webhook := router.Group(config.EndpointPrefix + "snap")
	webhook.POST("/payment/callback", snapController.HandleCallback)

	//User Standard Auth Routes (scope: public)
	auth := router.Group(config.EndpointPrefix + "user")
	auth.POST("/refresh", authController.RefreshToken)
	auth.POST("/login", authController.Login)
	auth.Use(gateMiddleware.HandleAuthAccess).POST("/register_email", authController.RegisterByEmail)
	auth.Use(gateMiddleware.HandleAuthAccess).POST("/otp", authController.VerifyOtp)

	//Post Ticketing (scope: public and admin in some cases)
	qr := router.Group(config.EndpointPrefix).Use(qrMiddleware.HandleScanQr)
	qr.GET("seat/:link", seatController.InfoByLink)

	//Pre Ticketing (scope: public)
	seatMap := router.Group(config.EndpointPrefix).Use(gateMiddleware.HandleTransactionAccess)
	seatMap.GET("seat_map", reservationController.GetSeatsInfo)

	//User data (scope: buyer user)
	user := router.Group(config.EndpointPrefix + "user").Use(userMiddleware.UserAccess)
	user.POST("/logout", authController.Logout)
	user.GET("/profile", userController.CurrentUser)
	user.GET("/tickets", userController.ShowMyTickets)
	user.PATCH("/profile", userController.UpdateInfo)

	//Ticketing routes (scope: buyer user)
	userTicketing := router.Group(config.EndpointPrefix).Use(gateMiddleware.HandleTransactionAccess).Use(userMiddleware.UserAccess)
	userTicketing.POST("seat_map", reservationController.ReserveSeats)
	userTicketing.GET("checkout", txController.GetLatestTransactionDetails)
	userTicketing.DELETE("checkout", txController.DeleteLatestTransaction)
	userTicketing.POST("checkout", txController.InitiateTransaction)

	//Admin Routes (scope: admin user)
	admin := router.Group(config.EndpointPrefix + "admin").Use(adminMiddleware.AdminAccess)
	admin.PUT("/seat/:link", seatController.UpdateByLink)

	admin.POST("/open_the_gate", gateController.OpenGate)
	admin.POST("/close_the_gate", gateController.CloseGate)

	admin.POST("/open_the_auth", gateController.OpenAuth)
	admin.POST("/close_the_auth", gateController.CloseAuth)

	admin.POST("/set_to_production", gateController.SetToProduction)
	admin.POST("/set_to_sandbox", gateController.SetToSandbox)

	admin.PATCH("/qr_scan_behaviour", gateController.UpdateQrScanBehaviour)
	admin.GET("/get_app_config", gateController.GetAppConfig)
	admin.GET("/seats", seatController.AllDetails)

	admin.GET("/healthAdmin", homeController.HealthCheck)

	//public.POST("/user/register", authController.Register) //This route is no longer needed for current GMCO's ticketing case,
	//public.POST("/user/sign_in", authController.SignIn) //but the code implementation in the controller is still remain in case of future use

	return router
}
