package controller

import (
	"errors"
	"fmt"
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
	config *config.AppConfig
	txDb   *gorm.DB
	log    *util.LogUtil

	reservationService *service.ReservationService
	txService          *service.TransactionService
	seatService        *service.SeatService
}

func NewReservationController(config *config.AppConfig, txDb *gorm.DB, log *util.LogUtil, reservationService *service.ReservationService, txService *service.TransactionService, seatService *service.SeatService) *ReservationController {
	return &ReservationController{config: config, txDb: txDb, log: log, reservationService: reservationService, txService: txService, seatService: seatService}
}

func (r *ReservationController) GetSeatsInfo(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion

	seats, err := r.seatService.GetAllSeats() //get all seats from db
	if err != nil {
		r.log.ControllerResponseLog(err, "ReservationController@GetSeatsInfo", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}

	seatsResponse := make([]validation.ReservationResponse, len(seats), len(seats)) //create response object
	for _, seat := range seats {
		seatsResponse[seat.SeatId-1].SeatId = seat.SeatId
		seatsResponse[seat.SeatId-1].Name = seat.Name
		seatsResponse[seat.SeatId-1].Status = seat.Status
		seatsResponse[seat.SeatId-1].Price = seat.Price
	}

	mySeats, _ := r.txService.SeatsBelongsToUser(accessDetails.UserId) //overwrite the response object for this user
	for _, mySeat := range mySeats {                                   //populate the response object
		seatsResponse[mySeat.SeatId-1].Status = mySeat.Status
	}

	for _, seat := range seats { //overwrite the response with timestamp logic
		if time.Now().After(seat.UpdatedAt.Add(r.config.TransactionMinute)) {
			seatsResponse[seat.SeatId-1].Status = "available"
		}
	}

	c.JSON(http.StatusOK, gin.H{ //return success
		"message": "success",
		"data":    seatsResponse,
		"count":   len(seatsResponse),
	})
	return
}

func (r *ReservationController) ReserveSeats(c *gin.Context) {
	contextData, _ := c.Get("accessDetails")              //get the details about the current user that make request from the context passed by user middleware
	accessDetails, _ := contextData.(*util.AccessDetails) //type assertion

	var inputData validation.ReservationRequest //get the seats data in request body
	if err := c.ShouldBindJSON(&inputData); err != nil {
		r.log.ControllerResponseLog(err, "ReservationController@ReserveSeats", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}

	if err := r.reservationService.CheckUserSeatCount(inputData.SeatIds, accessDetails.UserId); err != nil { //check user seat limit
		r.log.ControllerResponseLog(err, "ReservationController@ReserveSeats", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusForbidden, "error when processing the request data", err.Error())
		return
	}

	txn := r.txDb.Begin() //START DATABASE TRANSACTION
	if txn.Error != nil {
		fmt.Print(txn.Error)
	}

	for _, seatId := range inputData.SeatIds { //check eligibility for each chair in request
		if err := r.seatService.IsOwnedTxn(txn, seatId, accessDetails.UserId); err != nil {
			txn.Rollback() //ABORT DATABASE TRANSACTION
			err = errors.New(err.Error() + " | conflict on this seat. seat_id: " + strconv.Itoa(int(seatId)))
			r.log.ControllerResponseLog(err, "ReservationController@ReserveSeats", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
			util.GinResponseError(c, http.StatusConflict, "conflict when processing the request data", err.Error())
			return
		}
	}

	for _, seatId := range inputData.SeatIds { //update seat availability
		if err := r.seatService.UpdateStatusTxn(txn, seatId, "reserved"); err != nil {
			txn.Rollback() //ABORT DATABASE TRANSACTION
			r.log.ControllerResponseLog(err, "ReservationController@ReserveSeats", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
			util.GinResponseError(c, http.StatusConflict, "error when processing the request data", err.Error())
			return
		}
	}

	txcErr := txn.Commit().Error //COMMIT DATABASE TRANSACTION
	if txcErr != nil {
		fmt.Print(txcErr)
	}

	if err := r.txService.CreateTx(accessDetails.UserId, inputData.SeatIds); err != nil { //store reservation to txDb table
		r.log.ControllerResponseLog(err, "ReservationController@ReserveSeats", c.ClientIP(), contextData.(*util.AccessDetails).UserId)
		util.GinResponseError(c, http.StatusConflict, "error when processing the request data", err.Error())
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "information updated successfully",
	})
	return
}
