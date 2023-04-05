package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/app/validation"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *service.UserService
	tokenUtil   *util.TokenUtil
	config      *config.AppConfig
}

func NewUserController(userService *service.UserService, tokenUtil *util.TokenUtil, config *config.AppConfig) *UserController {
	return &UserController{userService: userService, tokenUtil: tokenUtil, config: config}
}

func (u *UserController) UpdateInfo(c *gin.Context) {
	contextData, isExist := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   "cannot get access details",
		})
		return
	}
	accessDetails, _ := contextData.(*util.AccessDetails)

	var inputData validation.RegisterValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	newUser := model.User{ //update the new user data
		Name:  inputData.Name,
		Email: inputData.Email,
		Phone: inputData.Phone,
	}
	affectedRows, err := u.userService.UpdateById(accessDetails.UserId, &newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "users info updated successfully",
		"affected_records": affectedRows,
	})

}
