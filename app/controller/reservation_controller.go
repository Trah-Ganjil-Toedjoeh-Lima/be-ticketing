package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReservationController struct {
	resSvc *service.ReservationService
}

func NewReservationController(resSvc *service.ReservationService) *ReservationController {
	return &ReservationController{resSvc: resSvc}
}

func (r *ReservationController) GetSeatsInfo(c *gin.Context) {
	seats, err := r.resSvc.GetAllSeats()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   seats,
	})
	return
}
