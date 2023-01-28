package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"strconv"
)

type TransactionService struct {
	txRepo   *repository.TransactionRepository
	userRepo *repository.UserRepository
	seatRepo *repository.SeatRepository
}

func NewTransactionService(txRepo *repository.TransactionRepository, userRepo *repository.UserRepository, seatRepo *repository.SeatRepository) *TransactionService {
	return &TransactionService{txRepo: txRepo, userRepo: userRepo, seatRepo: seatRepo}
}

func (s *TransactionService) CreateTx(userId uint64, seatIds []uint) error {
	txId := uuid.New().String()
	for _, seatId := range seatIds {
		//create tx for each seat
		newTx := model.Transaction{
			OrderId:      txId,
			UserId:       userId,
			SeatId:       seatId,
			Vendor:       "#",
			Confirmation: "reserved",
		}
		//delete previous failed reservation
		s.txRepo.SoftDeleteTransaction(seatId, userId)
		//save transaction
		if result := s.txRepo.InsertOne(&newTx); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (s *TransactionService) SeatsBelongsToUserId(userId uint64) ([]model.Seat, error) {
	var transactions []model.Transaction
	var seats []model.Seat
	if result := s.txRepo.GetLastTxByUserId(&transactions, userId); result.RowsAffected < 1 {
		return seats, errors.New("user belum melakukan pemesanan/transaksi")
	}

	for _, tx := range transactions {
		var seat model.Seat
		s.seatRepo.GetSeatById(&seat, tx.SeatId)
		if tx.Confirmation == "reserved" {
			seat.Status = "reserved_by_me"
		}
		if tx.Confirmation == "settlement" {
			seat.Status = "purchased_by_me"
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

func (s *TransactionService) GetUserTransactionDetails(userId uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetUserTransactionDetails(&transactions, userId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil

}

func (s *TransactionService) PrepareTransactionData(userId uint64) snap.Request {
	txDetails, _ := s.GetUserTransactionDetails(userId)
	var grossAmt int64
	var itemDetails []midtrans.ItemDetails

	customerDetails := midtrans.CustomerDetails{
		FName: txDetails[0].User.Name,
		LName: "",
		Email: txDetails[0].User.Email,
		Phone: txDetails[0].User.Phone,
	}

	for _, tx := range txDetails {
		grossAmt += int64(tx.Seat.Price)
		itemDetail := midtrans.ItemDetails{
			ID:    strconv.FormatUint(uint64(tx.SeatId), 10),
			Price: int64(tx.Seat.Price),
			Qty:   1,
			Name:  tx.Seat.Name,
		}
		itemDetails = append(itemDetails, itemDetail)
	}

	var snapRequest snap.Request = snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  txDetails[0].OrderId,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &customerDetails,
		Items:          &itemDetails,
	}
	return snapRequest
}

func (s *TransactionService) GetTxByOrderId(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetTxByOrderId(&transactions, orderId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) UpdatePaymentStatus(orderId, vendor, confirmation string) error {
	if result := s.txRepo.UpdatePaymentStatus(orderId, vendor, confirmation); result.Error != nil {
		return result.Error
	}
	return nil
}
