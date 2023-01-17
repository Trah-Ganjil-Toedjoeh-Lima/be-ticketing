package app

import (
	"github.com/frchandra/gmcgo/app/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

//type Router struct {}

//func NewRouter(map[string]interface{}) //TODO: hoarusnya map of controllers

func NewRouter(userController *controller.UserController) *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
		})
	})
	return router
}
