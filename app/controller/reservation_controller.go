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
	"reflect"
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
	userService        *service.UserService
	tokenUtil          *util.TokenUtil
}

func NewReservationController(config *config.AppConfig, txDb *gorm.DB, log *util.LogUtil, reservationService *service.ReservationService, txService *service.TransactionService, seatService *service.SeatService, userService *service.UserService, tokenUtil *util.TokenUtil) *ReservationController {
	return &ReservationController{config: config, txDb: txDb, log: log, reservationService: reservationService, txService: txService, seatService: seatService, userService: userService, tokenUtil: tokenUtil}
}

// GetSeatsInfo GET /seat_map
func (r *ReservationController) GetSeatsInfo(c *gin.Context) {
	seats, err := r.seatService.GetAllSeats() //get all seats from db
	if err != nil {
		r.log.ControllerResponseLog(err, "ReservationController@GetSeatsInfo", c.ClientIP(), 0)
		util.GinResponseError(c, http.StatusNotFound, "something went wrong", "error when getting the data")
		return
	}

	seatsResponse := make([]validation.ReservationResponse, len(seats), len(seats)) //create response object
	for _, seat := range seats {
		seatsResponse[seat.SeatId-1].SeatId = seat.SeatId
		seatsResponse[seat.SeatId-1].Name = seat.Name
		if seat.Status != "purchased" && time.Now().After(seat.UpdatedAt.Add(r.config.TransactionMinute)) { //overwrite the response with timestamp logic
			seat.Status = "available"
		}
		seatsResponse[seat.SeatId-1].Status = seat.Status
		seatsResponse[seat.SeatId-1].Price = seat.Price
		seatsResponse[seat.SeatId-1].Row = seat.Row
		seatsResponse[seat.SeatId-1].Column = seat.Column
	}

	accessDetails, tokenEmptyError := r.tokenUtil.GetValidatedAccess(c) //get the user data from the token in the request header
	if tokenEmptyError == nil {                                         //if credentials found (user is logged in) and token is not expired
		tokenExpiredError := r.tokenUtil.FetchAuthn(accessDetails.AccessUuid) //check if token exist in the token storage (Check if the token is expired)
		if tokenExpiredError == nil {
			mySeats, _ := r.txService.SeatsBelongsToUser(accessDetails.UserId) //overwrite the response object for this user
			for _, mySeat := range mySeats {                                   //populate the response object
				if seatsResponse[mySeat.SeatId-1].Status != "available" { //only overwrite the seat status if it was not overwritten previously by timestamp logic
					seatsResponse[mySeat.SeatId-1].Status = mySeat.Status
				}
			}
		}

	}

	c.JSON(http.StatusOK, gin.H{ //return success
		"message": "success",
		"data":    seatsResponse,
		"count":   len(seatsResponse),
	})
	return
}

// ReserveSeats POST /seat_map
func (r *ReservationController) ReserveSeats(c *gin.Context) {
	contextData, ok := c.Get("accessDetails") //get the details about the current user that make request from the context passed by user middleware
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": "cannot get access details"})
		return
	}
	accessDetails, ok := contextData.(*util.AccessDetails) //type assertion
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": "cannot processing access details"})
		return
	}

	user, err := r.userService.GetById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": err.Error()})
		return
	}
	if reflect.ValueOf(user.Name).IsZero() || reflect.ValueOf(user.Phone).IsZero() { //check if user has fill their identity data
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "fail",
			"error":   "you are not authorized, please fill your name or phone number data",
		})
		return
	}

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
