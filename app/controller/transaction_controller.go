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
	userSeats, _ := t.txService.SeatsBelongsToUserId(accessDetails.UserId)
	user, _ := t.userService.GetById(accessDetails.UserId)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   user,
		"seats":  userSeats,
	})
	return
}
