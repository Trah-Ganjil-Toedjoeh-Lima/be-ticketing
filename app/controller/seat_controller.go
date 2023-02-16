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
	transactions, err := s.txService.GetAllWithDetails()
	if err != nil {
		s.log.Log.WithField("occurrence", "SeatsController@InfoByLink").Error(err)
		util.GinResponseError(c, http.StatusNotFound, "request fail", "error when processing the request data")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    transactions,
	})
	return
}

func (s *SeatController) InfoByLink(c *gin.Context) {
	link := c.Param("link")
	seatDetails, err := s.txService.GetDetailsByLink(link)
	if err != nil {
		s.log.Log.WithField("occurrence", "SeatsController@InfoByLink").Error(err)
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

func (s *SeatController) DetailsByLink(c *gin.Context) {
	link := c.Param("link")
	seatDetails, err := s.txService.GetDetailsByLink(link)
	if err != nil {
		s.log.Log.WithField("occurrence", "SeatsController@InfoByLink").Error(err)
		util.GinResponseError(c, http.StatusNotFound, "request fail", "error when processing the request data")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    seatDetails,
	})
	return
}

func (s *SeatController) UpdateByLink(c *gin.Context) {
	link := c.Param("link")
	var inputData map[string]string //get the seats data in request body
	if err := c.ShouldBindJSON(&inputData); err != nil {
		s.log.BasicLog(err, "SeatController@UpdateByLink")
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}
	postSaleStatus := inputData["post_sale_status"]
	if err := s.seatService.UpdatePostSaleStatus(link, postSaleStatus); err != nil {
		s.log.BasicLog(err, "SeatController@UpdateByLink")
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    postSaleStatus,
	})
}

func (s *SeatController) UpdateToAttended(c *gin.Context) {
	link := c.Param("link")
	if err := s.seatService.UpdatePostSaleStatus(link, "attended"); err != nil {
		s.log.BasicLog(err, "SeatController@UpdateByLink")
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "attended",
	})
}

func (s *SeatController) UpdateToExchanged(c *gin.Context) {
	link := c.Param("link")
	if err := s.seatService.UpdatePostSaleStatus(link, "exchanged"); err != nil {
		s.log.BasicLog(err, "SeatController@UpdateByLink")
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    "exchanged",
	})
}
