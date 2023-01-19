package app

import (
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/frchandra/gmcgo/app/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userController *controller.UserController,
	userMiddleware *middleware.UserMiddleware,
) *gin.Engine {
	router := gin.Default()

	public := router.Group("/api/v1")

	public.POST("/user/register", userController.Register)
	public.POST("/user/login", userController.Login)
	public.POST("/user/refresh", userController.RefreshToken)

	user := router.Group("/api/v1").Use(userMiddleware.HandleUserAccess)
	user.POST("/user/logout", userController.Logout)
	user.GET("/user", userController.CurrentUser)

	return router
}
