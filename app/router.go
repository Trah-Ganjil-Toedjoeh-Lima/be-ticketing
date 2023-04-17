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

	//Health Check
	router.GET("/", homeController.HealthCheck)
	router.GET("/api/v1", homeController.HealthCheck)

	//Midtrans Webhook
	webhook := router.Group("api/v1")
	webhook.POST("/snap/payment/callback", snapController.HandleCallback)

	//Public User Standard Auth Routes
	public := router.Group("/api/v1")
	public.POST("/user/refresh", authController.RefreshToken)
	public.POST("/user/register_email", authController.RegisterByEmail)
	public.POST("/user/otp", authController.VerifyOtp)
	public.POST("/user/login", authController.Login)

	//Public Post Ticketing
	public.Use(qrMiddleware.HandleScanQr).GET("/seat/:link", seatController.InfoByLink)

	//Public Pre Ticketing
	public.Use(gateMiddleware.HandleAccess).GET("/seat_map", reservationController.GetSeatsInfo)

	//Logged-In User Routes
	user := router.Group("/api/v1").Use(userMiddleware.UserAccess)
	user.POST("/user/logout", authController.Logout)
	user.GET("/user", userController.CurrentUser)
	user.GET("/user/tickets", userController.ShowMyTickets)
	user.PATCH("/user", userController.UpdateInfo)

	user.Use(gateMiddleware.HandleAccess).POST("/seat_map", reservationController.ReserveSeats) //Ticketing routes
	user.Use(gateMiddleware.HandleAccess).GET("/checkout", txController.GetLatestTransactionDetails)
	user.Use(gateMiddleware.HandleAccess).POST("/checkout", txController.InitiateTransaction)

	//Admin Routes
	admin := router.Group("/api/v1").Use(adminMiddleware.AdminAccess)
	admin.GET("/admin/seat/:link", seatController.DetailsByLink)
	admin.PUT("/admin/seat/:link", seatController.UpdateByLink)
	admin.GET("/admin/seat/:link/:status", seatController.UpdateToStatus)

	admin.POST("/admin/open_the_gate", gateController.OpenGate)
	admin.POST("/admin/close_the_gate", gateController.CloseGate)

	admin.POST("/admin/set_to_production", gateController.SetToProduction)
	admin.POST("/admin/set_to_sandbox", gateController.SetToSandbox)

	admin.PATCH("/admin/qr_scan_behaviour", gateController.UpdateQrScanBehaviour)
	admin.GET("/admin/get_app_config", gateController.GetAppConfig)
	admin.GET("/admin/seats", seatController.AllDetails)

	//public.POST("/user/register", authController.Register) //This route is no longer needed for current GMCO's ticketing case,
	//public.POST("/user/sign_in", authController.SignIn) //but the code implementation in the controller is still remain in case of future use

	return router
}
