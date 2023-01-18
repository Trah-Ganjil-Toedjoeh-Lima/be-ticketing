package controller

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util/token"
	"github.com/frchandra/gmcgo/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userSercive *service.UserService
}

func NewUserController(userSercive *service.UserService) *UserController {
	return &UserController{userSercive: userSercive}
}

func (this *UserController) Register(c *gin.Context) {
	var userData validation.RegisterValidation
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	newUser := model.User{
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}
	rowsAffected, err := this.userSercive.InsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"rows_affected": rowsAffected,
	})
	return
}

func (this *UserController) Login(c *gin.Context) {
	var userData validation.LoginValidation
	var oldUser model.User
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	if userData.Name == "" {
		oldUser = model.User{
			Email: userData.Email,
			Phone: userData.Phone,
		}
	} else {
		oldUser = model.User{
			Name:  userData.Name,
			Phone: userData.Phone,
		}
	}
	token, err := this.userSercive.ValidateLogin(&oldUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", token, 3600, "/", "127.0.0.1", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return
}

func (this *UserController) CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := this.userSercive.GetById(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
	return
}
