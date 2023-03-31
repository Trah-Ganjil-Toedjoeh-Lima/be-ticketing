package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SnapController struct {
	snapService *service.SnapService
	snapUtil    *util.SnapUtil
	txService   *service.TransactionService
	log         *util.LogUtil
}

func NewSnapController(snapService *service.SnapService, snapUtil *util.SnapUtil, txService *service.TransactionService, log *util.LogUtil) *SnapController {
	return &SnapController{snapService: snapService, snapUtil: snapUtil, txService: txService, log: log}
}

func (s *SnapController) HandleCallback(c *gin.Context) {
	message := make(map[string]interface{})
	if err := c.ShouldBindJSON(&message); err != nil { //bind json from request body
		c.Status(http.StatusBadRequest)
		return
	}
	if err := s.snapUtil.CheckSignature(message); err != nil { //check the authenticity of the signature key from the json data
		c.Status(http.StatusBadRequest)
		return
	}
	txStatus, _ := message["transaction_status"].(string) //handle according to the "transaction_status" field from the json data
	if txStatus == "pending" {
		if err := s.snapService.HandlePending(message); err != nil {
			c.Status(http.StatusNotFound)
			s.log.BasicLog(err, "SnapController@HandleCallback@HandlePending")
			return
		}
		go func() {
			if err := s.snapService.SendInfoEmail(s.snapService.PrepareTxDetailsByMsg(message)); err != nil {
				s.log.BasicLog(err, "SnapController@HandleCallback@HandlePending@SendInfoEmail")
			}
		}()
	} else if txStatus == "settlement" {
		if err := s.snapService.HandleSettlement(message); err != nil {
			c.Status(http.StatusNotFound)
			s.log.BasicLog(err, "SnapController@HandleCallback@HandleSettlement")
			return
		}
		go func() {
			if err := s.snapService.SendTicketEmail(s.snapService.PrepareTxDetailsByMsg(message)); err != nil {
				s.log.BasicLog(err, "SnapController@HandleCallback@HandlePending@SendInfoEmail")
			}
		}()
	} else if txStatus == "expire" || txStatus == "cancel" || txStatus == "deny" {
		if err := s.snapService.HandleFailure(message); err != nil {
			s.log.BasicLog(err, "SnapController@HandleFailure@HandleSettlement")
			return
		}
	}
	c.Status(http.StatusOK)
	return

}
