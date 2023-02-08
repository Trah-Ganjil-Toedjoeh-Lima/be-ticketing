package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionController struct {
	txService   *service.TransactionService
	userService *service.UserService
	snapUtil    *util.SnapUtil
	log         *util.LogUtil
}

func NewTransactionController(txService *service.TransactionService, userService *service.UserService, snapUtil *util.SnapUtil, log *util.LogUtil) *TransactionController {
	return &TransactionController{txService: txService, userService: userService, snapUtil: snapUtil, log: log}
}

func (t *TransactionController) GetTransactionDetails(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion
	txDetails, err := t.txService.GetDetailsByUser(accessDetails.UserId)
	if err != nil {
		t.log.ControllerResponseLog(err, "TransactionController@GetTransactionDetails", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}

	var seatResponses []validation.SeatResponse //transform data
	for _, tx := range txDetails {
		seatResponse := validation.SeatResponse{Name: tx.Seat.Name, Price: tx.Seat.Price}
		seatResponses = append(seatResponses, seatResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"seats":      seatResponses,
			"user_name":  txDetails[0].User.Name,
			"user_email": txDetails[0].User.Email,
			"user_phone": txDetails[0].User.Phone,
		},
	})
	return
}

func (t *TransactionController) InitiateTransaction(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")                                     //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails)                        //type assertion
	snapRequest, err := t.txService.PrepareTransactionData(accessDetails.UserId) //prepare snap request
	if err != nil {
		t.log.ControllerResponseLog(err, "TransactionController@InitiateTransaction", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}
	response, midtransErr := t.snapUtil.CreateTransaction(&snapRequest) //send request to midtrans
	if midtransErr != nil {
		t.log.ControllerResponseLog(err, "TransactionController@InitiateTransaction", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		t.log.Log.
			WithField("snap_request", snapRequest).
			WithField("snap_response", response).
			Error("error when sending data to midtrans")
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"snap_request":  snapRequest,
		"snap_response": response,
	})
	return
}
