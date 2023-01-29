package controller

import (
	"fmt"
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
	//bind json
	if err := c.ShouldBindJSON(&message); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//check signature key
	if err := s.snapUtil.CheckSignature(message); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	//handle according tx status
	txStatus := message["transaction_status"].(string)
	if txStatus == "pending" {
		if err := s.snapService.HandlePending(message); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		go func() {
			if err := s.snapService.SendInfoEmail(s.snapService.PrepareTxDetailsByMsg(message)); err != nil {
				fmt.Println(err.Error())
			}
		}()
	} else if txStatus == "settlement" {
		if err := s.snapService.HandleSettlement(message); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		go func() {
			if err := s.snapService.SendTicketEmail(s.snapService.PrepareTxDetailsByMsg(message)); err != nil {
				fmt.Println(err.Error())
			}
		}()
	} else if txStatus == "expire" || txStatus == "failure" {
		if err := s.snapService.HandleFailure(message); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
	}

	c.Status(http.StatusOK)
	return

}
