package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/validation"
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
	seatsResponse := make([]validation.SeatResponse, len(seats), len(seats))
	for i, seat := range seats { //TODO: add is_reserved field to consumed by FE, this field can be dynamic for each user
		seatsResponse[i].SeatId = seat.SeatId
		seatsResponse[i].Name = seat.Name
		seatsResponse[i].Status = seat.Status
		seatsResponse[i].Price = seat.Price
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   seatsResponse,
		"count":  len(seatsResponse),
	})
	return
}

// TODO: finish this
func (r *ReservationController) ReserveSeats(c *gin.Context) {
	var requestData validation.SeatResrvRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

}
