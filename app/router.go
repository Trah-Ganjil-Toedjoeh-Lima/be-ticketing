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

	v1 := router.Group("/api/v1")

	v1.POST("/user/register", userController.Register)
	v1.POST("/user/login", userController.Login)
	v1.Use(userMiddleware.HandleUserAccess).POST("/user/logout", userController.Logout)
	v1.Use(userMiddleware.HandleUserAccess).GET("/me", userController.CurrentUser)

	return router
}
