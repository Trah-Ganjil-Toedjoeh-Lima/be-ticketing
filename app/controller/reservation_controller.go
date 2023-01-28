package controller

import (
	"errors"
	"github.com/frchandra/gmcgo/app/service"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/app/validation"
	"github.com/frchandra/gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ReservationController struct {
	resSvc      *service.ReservationService
	userService *service.UserService
	txService   *service.TransactionService
	seatService *service.SeatService
	config      *config.AppConfig
}

func NewReservationController(resSvc *service.ReservationService, userService *service.UserService, txService *service.TransactionService, seatService *service.SeatService, config *config.AppConfig) *ReservationController {
	return &ReservationController{resSvc: resSvc, userService: userService, txService: txService, seatService: seatService, config: config}
}

func (r *ReservationController) GetSeatsInfo(c *gin.Context) {
	//TODO: ganti dengan logic timestamp seperti di laravel
	//get all seats from db
	seats, err := r.resSvc.GetAllSeats()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err,
		})
		return
	}
	//create response object
	seatsResponse := make([]validation.SeatResponse, len(seats), len(seats))
	for _, seat := range seats { //TODO: add is_reserved field to consumed by FE, this field can be dynamic for each user
		seatsResponse[seat.SeatId-1].SeatId = seat.SeatId
		seatsResponse[seat.SeatId-1].Name = seat.Name
		seatsResponse[seat.SeatId-1].Status = seat.Status
		seatsResponse[seat.SeatId-1].Price = seat.Price
	}
	//get the details about the current user that make request from the context passed by user middleware
	contextData, _ := c.Get("accessDetails")
	//type assertion
	accessDetails, _ := contextData.(*util.AccessDetails)
	//verify that the user is present in the db
	if _, err := r.userService.GetById(accessDetails.UserId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err,
		})
		return
	}
	//if user exist, overwrite the response object for this user
	mySeats, _ := r.txService.SeatsBelongsToUserId(accessDetails.UserId)
	for _, mySeat := range mySeats {
		seatsResponse[mySeat.SeatId-1].Status = mySeat.Status
	}
	//overwrite with timestamp logic
	for _, seat := range seats {
		if time.Now().After(seat.UpdatedAt.Add(r.config.TransactionMinute)) {
			seatsResponse[seat.SeatId-1].Status = "available"
		}
	}

	//return success
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   seatsResponse,
		"count":  len(seatsResponse),
	})
	return
}

func (r *ReservationController) ReserveSeats(c *gin.Context) {
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

	//verify that the user is exists in the db
	if _, err := r.userService.GetById(accessDetails.UserId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	//get the requested seats
	var inputData validation.SeatReservationRequest
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
				"data":   err.Error(),
			})
			return

		}
	}

	//check user seat limit
	if err := r.resSvc.CheckUserSeatCount(inputData.SeatIds, accessDetails.UserId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "success",
			"data":   err.Error(),
		})
		return
	}

	//store reservation to tx table
	if err := r.txService.CreateTx(accessDetails.UserId, inputData.SeatIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	//update seat availability
	for _, seatId := range inputData.SeatIds {
		if err := r.seatService.UpdateStatus(seatId, "reserved"); err != nil { //TODO: ganti dengan logic timestamp seperti di laravel => tidak usah, pake update_at saja, langsung otomatis
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "fail",
				"data":   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"status": "success",
		"data":   inputData.SeatIds,
		"ok":     "ok",
	})
	return
}
