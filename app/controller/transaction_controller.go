package controller

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/app/validation"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionController struct {
	txService   *service.TransactionService
	userService *service.UserService
	snapUtil    *util.SnapUtil
	log         *util.LogUtil
	appConfig   *config.AppConfig
}

func NewTransactionController(txService *service.TransactionService, userService *service.UserService, snapUtil *util.SnapUtil, log *util.LogUtil, appConfig *config.AppConfig) *TransactionController {
	return &TransactionController{txService: txService, userService: userService, snapUtil: snapUtil, log: log, appConfig: appConfig}
}

// GetLatestTransactionDetails GET /checkout
func (t *TransactionController) GetLatestTransactionDetails(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //get the transaction and user info details for this request. user id is obtained from the context passed by user_middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion
	txDetails, err := t.txService.GetDetailsByUserConfirmation(accessDetails.UserId, "reserved")
	if err != nil {
		t.log.ControllerResponseLog(err, "TransactionController@GetLatestTransactionDetails", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}

	if len(txDetails) < 1 {
		t.log.ControllerResponseLog(errors.New("cannot find transaction data for this user"), "TransactionController@GetLatestTransactionDetails", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "cannot find data", "cannot find transaction data for this user")
		return
	}

	var seatResponses []validation.BasicSeatResponse //transform data
	for _, tx := range txDetails {
		seatResponse := validation.BasicSeatResponse{Name: tx.Seat.Name, Price: tx.Seat.Price, Category: tx.Seat.Category}
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

// InitiateTransaction POST /checkout
func (t *TransactionController) InitiateTransaction(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")                                     //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails)                        //type assertion
	snapRequest, err := t.txService.PrepareTransactionData(accessDetails.UserId) //prepare snap request
	if err != nil {
		t.log.ControllerResponseLog(err, "TransactionController@InitiateTransaction", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", err.Error())
		return
	}
	response, midtransErr := t.snapUtil.CreateTransaction(&snapRequest) //send request to midtrans
	if midtransErr != nil {
		t.log.ControllerResponseLog(midtransErr, "TransactionController@InitiateTransaction", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		t.log.Logrus.
			WithField("snap_request", snapRequest).
			WithField("snap_response", response).
			Error("error when sending data to midtrans")
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}

	var midtransClientKey string
	if t.appConfig.MidtransIsProduction == false {
		midtransClientKey = t.appConfig.ClientKeySandbox
	} else {
		midtransClientKey = t.appConfig.ClientKeyProduction
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "success",
		"snap_request":        snapRequest,
		"snap_response":       response,
		"midtrans_client_key": midtransClientKey,
	})
	return
}
