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
	txService   *service.TransactionService
	config      *config.AppConfig
}

func NewUserController(userService *service.UserService, txService *service.TransactionService, config *config.AppConfig) *UserController {
	return &UserController{userService: userService, txService: txService, config: config}
}

// UpdateInfo PATCH /user
func (u *UserController) UpdateInfo(c *gin.Context) {
	contextData, isExist := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   "cannot get access details",
		})
		return
	}
	accessDetails, ok := contextData.(*util.AccessDetails)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": "cannot get access details"})
		return
	}

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
			"message":       "fail",
			"error":         "This email is already registered. Please use other email",
			"error_details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "users info updated successfully",
		"affected_records": affectedRows,
	})
}

// CurrentUser GET /user/profile
func (u *UserController) CurrentUser(c *gin.Context) {
	contextData, isExist := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{"message": "fail", "error": "cannot get access details"})
		return
	}
	accessDetails, ok := contextData.(*util.AccessDetails)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": "cannot get access details"})
		return
	}

	user, err := u.userService.GetById(accessDetails.UserId) //get the user data given the user id from the token
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "fail", "error": err.Error()})
		return
	}

	var isAdmin bool
	if user.Name == u.config.AdminName && user.Phone == u.config.AdminPhone && user.Email == u.config.AdminEmail {
		isAdmin = true
	} else {
		isAdmin = false
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"data":     user,
		"is_admin": isAdmin,
	})
	return
}

// ShowMyTickets GET /user/tickets
func (u *UserController) ShowMyTickets(c *gin.Context) {
	contextData, isExist := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   "cannot get access details",
		})
		return
	}
	accessDetails, ok := contextData.(*util.AccessDetails)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": "cannot get access details"})
		return
	}

	transactions, err := u.txService.GetDetailsByUserConfirmation(accessDetails.UserId, "settlement")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "fail",
			"error":   err.Error(),
		})
		return
	}
	if len(transactions) < 1 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "fail",
			"error":   "this user does not have any transaction data yet",
		})
		return
	}

	var seatsResponse []validation.BasicSeatResponse
	var seatResponse validation.BasicSeatResponse
	for _, transaction := range transactions {
		seatResponse = validation.BasicSeatResponse{
			Name:     transaction.Seat.Name,
			Price:    transaction.Seat.Price,
			Category: transaction.Seat.Category,
		}
		seatsResponse = append(seatsResponse, seatResponse)
	}

	response := validation.UserTicketsResponse{
		Name:  transactions[0].User.Name,
		Phone: transactions[0].User.Phone,
		Email: transactions[0].User.Email,
		Seat:  seatsResponse,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response,
	})
	return

}
