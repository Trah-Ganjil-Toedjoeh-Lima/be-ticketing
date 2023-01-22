package controller

import (
	"errors"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/app/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReservationController struct {
	resSvc      *service.ReservationService
	userService *service.UserService
	txService   *service.TrsansactionService
}

func NewReservationController(resSvc *service.ReservationService, userService *service.UserService, txService *service.TrsansactionService) *ReservationController {
	return &ReservationController{resSvc: resSvc, userService: userService, txService: txService}
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
	//ambil informasi user
	//ambil informasi kursi yang akan dipesan
	//cek eligibility
	//simpan didatabase
	//return success with user data

	//get the details about the current user that make request from the context passed by user middleware
	contextData, isExist := c.Get("accessDetails")
	if isExist == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "cannot get access details",
		})
		return
	}
	//type assertion
	accessDetails, _ := contextData.(*util.AccessDetails)
	_, err := r.userService.GetById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	//get the requested seats
	var inputData validation.SeatResrvRequest
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	//check eligibility for each chair in request
	for _, seatId := range inputData.SeatIds {
		if err := r.resSvc.IsOwned(seatId, accessDetails.UserId); err != nil {
			err = errors.New(err.Error() + " [kursi ini :] " + strconv.Itoa(int(seatId)))
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "success",
				"data":   err,
			})
			break
		}
	}

	//store reservation to tx table

	c.JSON(http.StatusBadRequest, gin.H{
		"status": "success",
		"data":   inputData,
	})
	return
}
