package app

import (
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userController *controller.UserController,
) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	v1.POST("/user/register", userController.Register)
	v1.POST("/user/login", userController.Login)
	v1.GET("/me", userController.CurrentUser)

	return router
}
