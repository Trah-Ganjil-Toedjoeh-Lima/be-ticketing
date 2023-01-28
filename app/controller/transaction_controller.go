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
	snapUtil    *util.SnapUtil
}

func NewTransactionController(txService *service.TransactionService, userService *service.UserService, snapUtil *util.SnapUtil) *TransactionController {
	return &TransactionController{txService: txService, userService: userService, snapUtil: snapUtil}
}

func (t *TransactionController) GetTransactionDetails(c *gin.Context) { //TODO: make the response data more FE friendly
	//get the details about the current user that make request from the context passed by user middleware
	contextData, _ := c.Get("accessDetails")
	//type assertion
	accessDetails, _ := contextData.(*util.AccessDetails)
	txDetails, err := t.txService.GeTxDetailsByUser(accessDetails.UserId)
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
	//prepare snap request
	snapRequest := t.txService.PrepareTransactionData(accessDetails.UserId)
	//send request to midtrans
	response, err := t.snapUtil.CreateTransaction(&snapRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"err":    err.GetMessage(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"err":     response,
		"snapReq": snapRequest,
	})
	return
}
