package controller

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/app/validation"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type ReservationController struct {
	config             *config.AppConfig
	txDb               *gorm.DB
	reservationService *service.ReservationService
	userService        *service.UserService
	txService          *service.TransactionService
	seatService        *service.SeatService
}

func NewReservationController(config *config.AppConfig, txDb *gorm.DB, reservationService *service.ReservationService, userService *service.UserService, txService *service.TransactionService, seatService *service.SeatService) *ReservationController {
	return &ReservationController{config: config, txDb: txDb, reservationService: reservationService, userService: userService, txService: txService, seatService: seatService}
}

func (r *ReservationController) GetSeatsInfo(c *gin.Context) {
	seats, err := r.seatService.GetAllSeats() //get all seats from db
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	seatsResponse := make([]validation.SeatResponse, len(seats), len(seats)) //create response object
	for _, seat := range seats {
		seatsResponse[seat.SeatId-1].SeatId = seat.SeatId
		seatsResponse[seat.SeatId-1].Name = seat.Name
		seatsResponse[seat.SeatId-1].Status = seat.Status
		seatsResponse[seat.SeatId-1].Price = seat.Price
	}
	contextData, _ := c.Get("accessDetails")                             //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails)                //type assertion
	mySeats, _ := r.txService.IsSeatsBelongsToUser(accessDetails.UserId) //overwrite the response object for this user
	for _, mySeat := range mySeats {                                     //populate the response object
		seatsResponse[mySeat.SeatId-1].Status = mySeat.Status
	}
	for _, seat := range seats { //overwrite with timestamp logic
		if time.Now().After(seat.UpdatedAt.Add(r.config.TransactionMinute)) {
			seatsResponse[seat.SeatId-1].Status = "available"
		}
	}
	c.JSON(http.StatusOK, gin.H{ //return success
		"status": "success",
		"data":   seatsResponse,
		"count":  len(seatsResponse),
	})
	return
}

func (r *ReservationController) ReserveSeats(c *gin.Context) {
	contextData, isExist := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if isExist == false {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "cannot get access details",
		})
		return
	}
	accessDetails, _ := contextData.(*util.AccessDetails)                  //type assertion
	if _, err := r.userService.GetById(accessDetails.UserId); err != nil { //verify that the user is exists in the db
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	var inputData validation.SeatReservationRequest //get the seats data in request body
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}
	if err := r.reservationService.CheckUserSeatCount(inputData.SeatIds, accessDetails.UserId); err != nil { //check user seat limit
		c.JSON(http.StatusConflict, gin.H{
			"status": "success",
			"data":   err.Error(),
		})
		return
	}
	r.txDb.Begin()                             //START DATABASE TRANSACTION
	for _, seatId := range inputData.SeatIds { //check eligibility for each chair in request
		if err := r.seatService.IsOwned(seatId, accessDetails.UserId); err != nil {
			r.txDb.Rollback() //ABORT DATABASE TRANSACTION
			err = errors.New(err.Error() + " | conflict on this seat. seat_id: " + strconv.Itoa(int(seatId)))
			c.JSON(http.StatusConflict, gin.H{
				"status": "success",
				"data":   err.Error(),
			})
			return

		}
	}
	for _, seatId := range inputData.SeatIds { //update seat availability
		if err := r.seatService.UpdateStatus(seatId, "reserved"); err != nil {
			r.txDb.Rollback() //ABORT DATABASE TRANSACTION
			c.JSON(http.StatusConflict, gin.H{
				"status": "fail",
				"data":   err.Error(),
			})
			return
		}
	}
	r.txDb.Commit()                                                                       //COMMIT DATABASE TRANSACTION
	if err := r.txService.CreateTx(accessDetails.UserId, inputData.SeatIds); err != nil { //store reservation to txDb table
		c.JSON(http.StatusConflict, gin.H{
			"status": "fail",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   inputData.SeatIds,
		"ok":     "ok",
	})
	return
}
