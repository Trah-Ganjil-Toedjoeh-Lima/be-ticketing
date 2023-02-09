package app

import (
	"github.com/frchandra/ticketing-gmcgo/app/controller"
	"github.com/frchandra/ticketing-gmcgo/app/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userMiddleware *middleware.UserMiddleware,
	adminMiddleware *middleware.AdminMiddleware,
	gateMiddleware *middleware.GateMiddleware,

	userController *controller.UserController,
	reservationController *controller.ReservationController,
	txController *controller.TransactionController,
	snapController *controller.SnapController,
	gateController *controller.GateController,
	seatController *controller.SeatController,
) *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1").Use(gateMiddleware.HandleAccess)
	public.POST("/user/register", userController.Register)
	public.POST("/user/sign_in", userController.SignIn)
	public.POST("/user/login", userController.Login)
	public.POST("/user/refresh", userController.RefreshToken)

	public.GET("/seat/:link", seatController.InfoByLink)
	public.PUT("/seat/:link", seatController.UpdateByLink)

	webhook := router.Group("api/v1")
	webhook.POST("/snap/payment/callback", snapController.HandleCallback)

	user := router.Group("/api/v1").Use(gateMiddleware.HandleAccess).Use(userMiddleware.HandleUserAccess)
	user.POST("/user/logout", userController.Logout)
	user.GET("/user", userController.CurrentUser)
	user.Use(gateMiddleware.HandleAccess).GET("/seat_map", reservationController.GetSeatsInfo)
	user.Use(gateMiddleware.HandleAccess).POST("/seat_map", reservationController.ReserveSeats)
	user.GET("/checkout", txController.GetTransactionDetails)
	user.POST("/checkout", txController.InitiateTransaction)

	admin := router.Group("/api/v1").Use(adminMiddleware.HandleAdminAccess)
	admin.POST("/admin/open_the_gate", gateController.OpenGate)
	admin.POST("/admin/close_the_gate", gateController.CloseGate)
	admin.GET("/admin/get_app_config", gateController.GetAppConfig)

	return router
}
