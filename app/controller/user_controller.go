package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userSercive *service.UserService
}

func NewUserController(userSercive *service.UserService) *UserController {
	return &UserController{userSercive: userSercive}
}

func (this *UserController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})

	return
}
