package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SeatController struct {
	seatService *service.SeatService
	txService   *service.TransactionService
	log         *util.LogUtil
}

func NewSeatController(seatService *service.SeatService, txService *service.TransactionService, log *util.LogUtil) *SeatController {
	return &SeatController{seatService: seatService, txService: txService, log: log}
}

func (s *SeatController) AllDetails(c *gin.Context) {

}

func (s *SeatController) DetailsByLink(c *gin.Context) {
	link := c.Param("link")
	seatDetails, err := s.txService.GetDetailsByLink(link)
	if err != nil {
		s.log.Log.WithField("occurrence", "SeatsController@DetailsByLink").Error(err)
		util.GinResponseError(c, http.StatusNotFound, "request fail", "error when processing the request data")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"user_name":  seatDetails.User.Name,
			"user_email": seatDetails.User.Email,
			"user_phone": seatDetails.User.Phone,
			"seat_name":  seatDetails.Seat.Name,
		},
	})
	return

}

func (s *SeatController) UpdateByLink(c *gin.Context) {

}
