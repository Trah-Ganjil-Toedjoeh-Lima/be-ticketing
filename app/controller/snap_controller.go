package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SnapController struct {
	snapService *service.SnapService
	snapUtil    *util.SnapUtil
	txService   *service.TransactionService
}

func NewSnapController(snapService *service.SnapService, snapUtil *util.SnapUtil, txService *service.TransactionService) *SnapController {
	return &SnapController{snapService: snapService, snapUtil: snapUtil, txService: txService}
}

func (s *SnapController) HandleCallback(c *gin.Context) {
	message := make(map[string]interface{})
	if err := c.ShouldBindJSON(&message); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := s.snapUtil.CheckSignature(message); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	txStatus := message["transaction_status"].(string)
	if txStatus == "pending" {
		if err := s.snapService.HandlePending(message); err != nil {
			c.Status(http.StatusNotFound)
		}
		go func() {
			_ = s.snapService.SendInfoEmail(s.snapService.PrepareTxDetailsByMsg(message))
		}()
	} else if txStatus == "settlement" {
		if err := s.snapService.HandleSettlement(message); err != nil {
			c.Status(http.StatusNotFound)
		}
		go func() {
			_ = s.snapService.SendTicketEmail(s.snapService.PrepareTxDetailsByMsg(message))
		}()
	} else if txStatus == "expire" || txStatus == "failure" {
		if err := s.snapService.HandleFailure(message); err != nil {
			c.Status(http.StatusNotFound)
		}
	}

	c.Status(200)
	return

}
