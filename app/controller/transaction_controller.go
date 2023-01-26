package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionController struct {
	txService   *service.TransactionService
	userService *service.UserService
}

func NewTransactionController(txService *service.TransactionService, userService *service.UserService) *TransactionController {
	return &TransactionController{txService: txService, userService: userService}
}

func (t *TransactionController) GetTransactionDetails(c *gin.Context) {
	//get the details about the current user that make request from the context passed by user middleware
	contextData, _ := c.Get("accessDetails")
	//type assertion
	accessDetails, _ := contextData.(*util.AccessDetails)
	txDetails, err := t.txService.GetUserTransactionDetails(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"err":    err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   txDetails,
	})
	return
}

func (t *TransactionController) InitiateTransaction(c *gin.Context) {
	//get the details about the current user that make request from the context passed by user middleware
	contextData, _ := c.Get("accessDetails")
	//type assertion
	accessDetails, _ := contextData.(*util.AccessDetails)

	/*	//get seats reserved by user
		userSeats, err := t.txService.SeatsBelongsToUserId(accessDetails.UserId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "fail",
				"err":    err.Error(),
			})
			return
		}
		//get user details data
		user, err := t.userService.GetById(accessDetails.UserId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": "fail",
				"err":    err.Error(),
			})
			return
		}*/
	c.JSON(http.StatusOK, gin.H{
		"status": "fail",
		"err":    t.txService.GetUserTransactionDetails(accessDetails.UserId),
	})
	return
}
