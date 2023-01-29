package controller

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *service.UserService
	tokenUtil   *util.TokenUtil
}

func NewUserController(userService *service.UserService, tokenUtil *util.TokenUtil) *UserController {
	return &UserController{
		userService: userService,
		tokenUtil:   tokenUtil,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	//validate the input data
	var inputData validation.RegisterValidation
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//upsert the new user data
	newUser := model.User{
		Name:  inputData.Name,
		Email: inputData.Email,
		Phone: inputData.Phone,
	}
	_, err := uc.userService.GetOrInsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//generate token for this user
	token, err := uc.userService.GenerateToken(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//return success
	c.SetSameSite(http.SameSiteNoneMode)
	//c.SetCookie("token", token, 3600, "/", "127.0.0.1", false, true)
	//TODO: cookie?
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return
}

func (uc *UserController) SignIn(c *gin.Context) {
	//validate the input data
	var inputData validation.RegisterValidation
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//insert the new user data
	newUser := model.User{
		Name:  inputData.Name,
		Email: inputData.Email,
		Phone: inputData.Phone,
	}
	rowsAffected, err := uc.userService.InsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//return success
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"rows_affected": rowsAffected,
	})
	return
}

func (uc *UserController) Login(c *gin.Context) {
	var inputData validation.LoginValidation
	var userInput model.User
	//validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//choose between the given credential. Can be user's name or email
	if inputData.Name == "" {
		userInput = model.User{
			Email: inputData.Email,
			Phone: inputData.Phone,
		}
	} else {
		userInput = model.User{
			Name:  inputData.Name,
			Phone: inputData.Phone,
		}
	}
	//validate if user exist and credential is correct
	err := uc.userService.ValidateLogin(&userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//generate token for this user
	token, err := uc.userService.GenerateToken(&userInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	//return success
	c.SetSameSite(http.SameSiteNoneMode)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return
}

func (uc *UserController) CurrentUser(c *gin.Context) {
	//get the details about the current user that make request from the context passed by user middleware
	contextData, isExist := c.Get("accessDetails")
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "cannot get access details",
		})
		return
	}
	accessDetails, _ := contextData.(*util.AccessDetails)
	//get the user data given the user id from the token
	user, err := uc.userService.GetById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
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

func (uc *UserController) Logout(c *gin.Context) {
	accessDetails, err := uc.tokenUtil.GetValidatedAccess(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deleted, err := uc.tokenUtil.DeleteAuthn(accessDetails.AccessUuid)
	if err != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (uc *UserController) RefreshToken(c *gin.Context) {
	token, err := uc.tokenUtil.Refresh(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return

}
