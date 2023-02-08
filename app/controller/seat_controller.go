package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/gin-gonic/gin"
)

type SeatController struct {
	seatService *service.SeatService
}

func NewSeatController(seatService *service.SeatService) *SeatController {
	return &SeatController{seatService: seatService}
}

func (s *SeatController) AllDetails(c *gin.Context) {

}

func (s *SeatController) DetailsByLink(c *gin.Context) {
	link := c.Param("link")
	seatDetails := s.seatService.
}

func (s *SeatController) UpdateByLink(c *gin.Context) {

}


