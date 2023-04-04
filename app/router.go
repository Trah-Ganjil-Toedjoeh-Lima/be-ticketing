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

	userController *controller.UserController,
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

	//Public User Standard Auth Routes
	public := router.Group("/api/v1").Use(gateMiddleware.HandleAccess)
	public.POST("/user/register", userController.Register)
	public.POST("/user/sign_in", userController.SignIn)
	public.POST("/user/login", userController.Login)
	public.POST("/user/refresh", userController.RefreshToken)
	public.POST("user/register_email", userController.RegisterEmail)

	//Public Post Ticketing
	public.Use(qrMiddleware.HandleScanQr).GET("/seat/:link", seatController.InfoByLink)

	//Midtrans Webhook
	webhook := router.Group("api/v1")
	webhook.POST("/snap/payment/callback", snapController.HandleCallback)

	//Logged-In User Routes
	user := router.Group("/api/v1").Use(gateMiddleware.HandleAccess).Use(userMiddleware.HandleUserAccess)
	user.POST("/user/logout", userController.Logout)
	user.GET("/user", userController.CurrentUser)

	//Ticketing Routes
	user.Use(gateMiddleware.HandleAccess).GET("/seat_map", reservationController.GetSeatsInfo)
	user.Use(gateMiddleware.HandleAccess).POST("/seat_map", reservationController.ReserveSeats)
	user.GET("/checkout", txController.GetNewTransactionDetails)
	user.POST("/checkout", txController.InitiateTransaction)

	//Admin Routes
	admin := router.Group("/api/v1").Use(adminMiddleware.HandleAdminAccess)
	admin.GET("/admin/seat/:link", seatController.DetailsByLink)
	admin.PUT("/admin/seat/:link", seatController.UpdateByLink)
	admin.GET("/admin/seat/attended/:link", seatController.UpdateToAttended)
	admin.GET("/admin/seat/exchanged/:link", seatController.UpdateToExchanged)
	admin.POST("/admin/open_the_gate", gateController.OpenGate)
	admin.POST("/admin/close_the_gate", gateController.CloseGate)
	admin.PATCH("/admin/qr_scan_behaviour", gateController.UpdateQrScanBehaviour)
	admin.GET("/admin/get_app_config", gateController.GetAppConfig)

	return router
}
