package app

import (
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/frchandra/gmcgo/app/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userMiddleware *middleware.UserMiddleware,

	userController *controller.UserController,
	reservationController *controller.ReservationController,
	txController *controller.TransactionController,
) *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1")

	public.POST("/user/register", userController.Register)
	public.POST("/user/sign_in", userController.SignIn)
	public.POST("/user/login", userController.Login)
	public.POST("/user/refresh", userController.RefreshToken)

	user := router.Group("/api/v1").Use(userMiddleware.HandleUserAccess)
	user.POST("/user/logout", userController.Logout)
	user.GET("/user", userController.CurrentUser)
	user.GET("/seat_map", reservationController.GetSeatsInfo)
	user.POST("/seat_map", reservationController.ReserveSeats)
	user.GET("/checkout", txController.GetTransactionDetails)

	return router
}
