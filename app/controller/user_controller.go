package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/app/validation"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"net/http"
	"time"
)

type UserController struct {
	userService *service.UserService
	tokenUtil   *util.TokenUtil
	config      *config.AppConfig
}

func NewUserController(userService *service.UserService, tokenUtil *util.TokenUtil, config *config.AppConfig) *UserController {
	return &UserController{userService: userService, tokenUtil: tokenUtil, config: config}
}

func (u *UserController) RegisterEmail(c *gin.Context) {
	var inputData validation.RegisterEmailValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	newTotpSecret, err := totp.Generate(totp.GenerateOpts{ //generate newTotpSecret for this user
		Issuer:      "gmco",
		AccountName: inputData.Email,
		Period:      300,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	newUser := model.User{ // insert new user record
		Email: inputData.Email,
		Otp:   newTotpSecret.Secret(),
	}
	_, err = u.userService.GetOrInsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	newOtpToken, err := totp.GenerateCodeCustom(newTotpSecret.Secret(), time.Now(), totp.ValidateOpts{Period: 300})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{ //TODO: sent the otp token from email
		"message":    "success",
		"otp_token":  newTotpSecret.Secret(),
		"otp_secret": newOtpToken,
	})

}

func (u *UserController) Register(c *gin.Context) {
	var inputData validation.RegisterValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	newUser := model.User{ //upsert the new user data
		Name:  inputData.Name,
		Email: inputData.Email,
		Phone: inputData.Phone,
	}
	_, err := u.userService.GetOrInsertOne(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&newUser) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{ //return success
		"status": "success",
		"token":  token,
	})
	return
}

func (u *UserController) SignIn(c *gin.Context) {
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
	rowsAffected, err := u.userService.InsertOne(&newUser)
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

func (u *UserController) Login(c *gin.Context) {
	var inputData validation.LoginValidation
	var userInput model.User

	if err := c.ShouldBindJSON(&inputData); err != nil { //validate the input data
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	if inputData.Name == "" { //choose between the given credential. Can be user's name or email
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

	err := u.userService.ValidateLogin(&userInput) //validate if user exist and credential is correct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&userInput) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode) //return success
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
	return
}

func (u *UserController) CurrentUser(c *gin.Context) {
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
	user, err := u.userService.GetById(accessDetails.UserId)
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

func (u *UserController) Logout(c *gin.Context) {
	accessDetails, err := u.tokenUtil.GetValidatedAccess(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deleted, err := u.tokenUtil.DeleteAuthn(accessDetails.AccessUuid)
	if err != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
	return
}

func (u *UserController) RefreshToken(c *gin.Context) {
	token, err := u.tokenUtil.Refresh(c)
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
