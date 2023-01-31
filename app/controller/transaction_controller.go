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

func (t *TransactionController) GetTransactionDetails(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion
	txDetails, err := t.txService.GetTxDetailsByUser(accessDetails.UserId)
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
	contextData, _ := c.Get("accessDetails")                                     //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails)                        //type assertion
	snapRequest, err := t.txService.PrepareTransactionData(accessDetails.UserId) //prepare snap request
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	response, midtransErr := t.snapUtil.CreateTransaction(&snapRequest) //send request to midtrans
	if midtransErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "fail",
			"err":    midtransErr.GetMessage(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"err":     response,
		"snapReq": snapRequest,
	})
	return
}
