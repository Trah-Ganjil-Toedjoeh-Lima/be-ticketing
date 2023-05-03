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
	"reflect"
	"time"
)

type AuthController struct {
	userService *service.UserService
	tokenUtil   *util.TokenUtil
	log         *util.LogUtil
	emailUtil   *util.EmailUtil
	config      *config.AppConfig
}

func NewAuthController(userService *service.UserService, tokenUtil *util.TokenUtil, log *util.LogUtil, emailUtil *util.EmailUtil, config *config.AppConfig) *AuthController {
	return &AuthController{userService: userService, tokenUtil: tokenUtil, log: log, emailUtil: emailUtil, config: config}
}

// VerifyOtp POST /user/otp
func (u *AuthController) VerifyOtp(c *gin.Context) {
	var inputData validation.OtpVerification //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	user, err := u.userService.GetByEmail(inputData.Email) //get user from db using the email field
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": err.Error()})
		return
	}

	if valid, err := totp.ValidateCustom(inputData.Otp, user.TotpSecret, time.Now(), totp.ValidateOpts{Period: u.config.TotpPeriod, Algorithm: otp.AlgorithmSHA1, Digits: 6, Skew: 1}); err != nil { //verify otp
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	} else if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "otp is not valid",
		})
		return
	}

	token, err := u.userService.GenerateToken(&user) //generate access token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{ //return success
		"message": "success",
		"token":   token,
	})
	return

}

func (u *AuthController) RegisterByEmail(c *gin.Context) {
	var inputData validation.RegisterEmailValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	user, _ := u.userService.GetByEmail(inputData.Email) //check if user exist in the database

	if reflect.ValueOf(user).IsZero() { //if the user not found
		totpSecret, err := totp.Generate(totp.GenerateOpts{ //generate newTotpSecret for this user
			Issuer:      "gmco",
			AccountName: inputData.Email,
			Period:      u.config.TotpPeriod,
			Algorithm:   otp.AlgorithmSHA1,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail",
				"error":   err.Error(),
			})
			return
		}

		newUser := model.User{ // insert new user record
			Email:      inputData.Email,
			TotpSecret: totpSecret.Secret(),
		}
		_, err = u.userService.InsertOne(&newUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail",
				"error":   err.Error(),
			})
			return
		}

		totpToken, err := totp.GenerateCodeCustom(totpSecret.Secret(), time.Now(), totp.ValidateOpts{Period: u.config.TotpPeriod}) //generate totp token
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "fail",
				"error":   err.Error(),
			})
			return
		}

		go func() { //send the totp to the registered email
			data := map[string]any{"Totp": totpToken}
			if err = u.emailUtil.SendTotpEmail(data, inputData.Email); err != nil {
				u.log.BasicLog(err, "When sending totp")
				return
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"message":             "success",
			"totp_token":          totpToken,
			"is_new_registration": true,
		})
		return
	}

	totpToken, err := totp.GenerateCodeCustom(user.TotpSecret, time.Now(), totp.ValidateOpts{Period: u.config.TotpPeriod}) //generate totp token
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	go func() { //send the totp to the registered email
		data := map[string]any{"Totp": totpToken}
		if err = u.emailUtil.SendTotpEmail(data, inputData.Email); err != nil {
			u.log.BasicLog(err, "When sending totp")
			return
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message":             "success",
		"totp_token":          totpToken,
		"is_new_registration": false,
	})
}

func (u *AuthController) Register(c *gin.Context) {
	var inputData validation.RegisterValidation //validate the input data
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
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
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&newUser) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{ //return success
		"message": "success",
		"token":   token,
	})
	return
}

func (u *AuthController) SignIn(c *gin.Context) {
	//validate the input data
	var inputData validation.RegisterValidation
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
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
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}
	//return success
	c.JSON(http.StatusOK, gin.H{
		"message":       "success",
		"rows_affected": rowsAffected,
	})
	return
}

// Login POST /user/login
func (u *AuthController) Login(c *gin.Context) {
	var inputData validation.LoginValidation
	var userInput model.User

	if err := c.ShouldBindJSON(&inputData); err != nil { //validate the input data
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
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
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	token, err := u.userService.GenerateToken(&userInput) //generate token for this user
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode) //return success
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
	return
}

// Logout POST /user/logout
func (u *AuthController) Logout(c *gin.Context) {
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
		"message": "success",
	})
	return
}

// RefreshToken POST /user/refresh
func (u *AuthController) RefreshToken(c *gin.Context) {
	token, err := u.tokenUtil.Refresh(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
	return

}
