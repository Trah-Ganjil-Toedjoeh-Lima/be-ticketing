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

	apiv1 := "/api/v1"

	//Health Check (scope: public)
	router.GET("/", homeController.HealthCheck)
	router.GET(apiv1, homeController.HealthCheck)

	//Midtrans Webhook (scope: public)
	webhook := router.Group(apiv1)
	webhook.POST("/snap/payment/callback", snapController.HandleCallback)

	//User Standard Auth Routes (scope: public)
	auth := router.Group(apiv1 + "/user")
	auth.POST("/refresh", authController.RefreshToken)
	auth.POST("/register_email", authController.RegisterByEmail)
	auth.POST("/otp", authController.VerifyOtp)
	auth.POST("/login", authController.Login)

	//Post Ticketing (scope: public)
	qr := router.Group(apiv1).Use(qrMiddleware.HandleScanQr)
	qr.GET("/seat/:link", seatController.InfoByLink)

	//Pre Ticketing (scope: public)
	seatMap := router.Group(apiv1).Use(gateMiddleware.HandleAccess)
	seatMap.GET("/seat_map", reservationController.GetSeatsInfo)

	//User data (scope: buyer user)
	user := router.Group(apiv1 + "/user").Use(userMiddleware.UserAccess)
	user.POST("/logout", authController.Logout)
	user.GET("/", userController.CurrentUser)
	user.GET("/tickets", userController.ShowMyTickets)
	user.PATCH("/", userController.UpdateInfo)

	//Ticketing routes (scope: buyer user)
	userTicketing := router.Group(apiv1).Use(gateMiddleware.HandleAccess)
	userTicketing.POST("/seat_map", reservationController.ReserveSeats)
	userTicketing.GET("/checkout", txController.GetLatestTransactionDetails)
	userTicketing.POST("/checkout", txController.InitiateTransaction)

	//Admin Routes (scope: admin user)
	admin := router.Group(apiv1 + "/admin").Use(adminMiddleware.AdminAccess)
	admin.GET("/seat/:link", seatController.DetailsByLink)
	admin.PUT("/seat/:link", seatController.UpdateByLink)
	admin.GET("/seat/:link/:status", seatController.UpdateToStatus)

	admin.POST("/open_the_gate", gateController.OpenGate)
	admin.POST("/close_the_gate", gateController.CloseGate)

	admin.POST("/set_to_production", gateController.SetToProduction)
	admin.POST("/set_to_sandbox", gateController.SetToSandbox)

	admin.PATCH("/qr_scan_behaviour", gateController.UpdateQrScanBehaviour)
	admin.GET("/get_app_config", gateController.GetAppConfig)
	admin.GET("/seats", seatController.AllDetails)

	//public.POST("/user/register", authController.Register) //This route is no longer needed for current GMCO's ticketing case,
	//public.POST("/user/sign_in", authController.SignIn) //but the code implementation in the controller is still remain in case of future use

	return router
}
