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
	v1.POST("/user", userController.Register)

	return router
}
