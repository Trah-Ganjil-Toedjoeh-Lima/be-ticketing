package service

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/repository"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"strconv"
	"time"
)

type TransactionService struct {
	txRepo *repository.TransactionRepository
	config *config.AppConfig
}

func NewTransactionService(txRepo *repository.TransactionRepository, config *config.AppConfig) *TransactionService {
	return &TransactionService{txRepo: txRepo, config: config}
}

func (s *TransactionService) CreateTx(userId uint64, seatIds []uint) error {
	for _, seatId := range seatIds { //create tx for each seat
		newTx := model.Transaction{
			OrderId:      "",
			UserId:       userId,
			SeatId:       seatId,
			Vendor:       "no_vendor",
			Confirmation: "reserved",
		}
		s.txRepo.SoftDeleteBySeatUser(seatId, userId)                  //delete the previous failed reservation
		if result := s.txRepo.InsertOne(&newTx); result.Error != nil { //save the transaction
			return result.Error
		}
	}
	return nil
}

func (s *TransactionService) GetAllWithDetails() ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetAllWithDetails(&transactions); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) GetDetailsByLink(link string) (model.Transaction, error) {
	var transaction model.Transaction
	if result := s.txRepo.GetDetailsByLink(&transaction, link); result.Error != nil {
		return transaction, result.Error
	}
	return transaction, nil
}

func (s *TransactionService) GetBasicsByLink(link string) (model.Transaction, error) {
	var transaction model.Transaction
	if result := s.txRepo.GetBasicsByLink(&transaction, link); result.Error != nil {
		return transaction, result.Error
	}
	return transaction, nil
}

func (s *TransactionService) GetByUser(userId uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUser(&transactions, userId); result.Error != nil {
		return transactions, result.Error
	}
	transactions = s.CleanUpGhostTransaction(transactions)
	return transactions, nil
}

func (s *TransactionService) GetDetailsByUserConfirmation(userId uint64, confirmation string) ([]model.Transaction, error) {
	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUserConfirmation(&transactions, userId, confirmation); result.Error != nil {
		return transactions, result.Error
	}
	transactions = s.CleanUpGhostTransaction(transactions)
	return transactions, nil
}

func (s *TransactionService) GetByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetByOrder(&transactions, orderId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) GetDetailsByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetDetailsByOrder(&transactions, orderId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) CleanUpGhostTransaction(transactions []model.Transaction) []model.Transaction {
	var newTransaction []model.Transaction //cek apakah ada transaksi ngambang, jika ada buang dari slice dan update db
	for _, tx := range transactions {
		if time.Now().After(tx.CreatedAt.Add(s.config.TransactionMinute)) && tx.Confirmation != "settlement" { //if tx created_at + 15 < time now  => berarti transaction ngambang
			s.txRepo.UpdatePaymentStatusById(tx.TransactionId, "not_continued") //update database
			s.txRepo.SoftDeleteBySeatUser(tx.Seat.SeatId, tx.User.UserId)
		} else {
			newTransaction = append(newTransaction, tx)
		}
	}
	return newTransaction //transaksi bersih
}

func (s *TransactionService) SeatsBelongsToUser(userId uint64) ([]model.Seat, error) {
	var seats []model.Seat

	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUser(&transactions, userId); result.Error != nil {
		return seats, result.Error
	}
	if transactions = s.CleanUpGhostTransaction(transactions); len(transactions) < 1 {
		return seats, errors.New("this user does not have any transaction")
	}
	for _, tx := range transactions {
		if tx.Confirmation == "reserved" {
			tx.Seat.Status = "reserved_by_me"
		}
		if tx.Confirmation == "settlement" {
			tx.Seat.Status = "purchased_by_me"
		}
		seats = append(seats, tx.Seat)
	}
	return seats, nil
}

func (s *TransactionService) PrepareTransactionData(userId uint64) (snap.Request, error) {
	var txDetails []model.Transaction
	s.txRepo.GetDetailsByUserConfirmation(&txDetails, userId, "reserved")     //get user's transaction
	if txDetails = s.CleanUpGhostTransaction(txDetails); len(txDetails) < 1 { //clean up 'ghost' transaction that may be created by this user
		return snap.Request{}, errors.New("cannot find any transaction for this user")
	}
	orderId := uuid.New().String() //create order_id for the new midtrans transaction

	customerDetails := midtrans.CustomerDetails{ //populate the midtrans request with the customer detail
		FName: txDetails[0].User.Name,
		LName: "",
		Email: txDetails[0].User.Email,
		Phone: txDetails[0].User.Phone,
	}
	var grossAmt int64 //populate the item detail
	var itemDetails []midtrans.ItemDetails
	for _, tx := range txDetails {
		s.txRepo.UpdateOrderIdById(tx.TransactionId, orderId) //update order_id of this transaction in the database

		grossAmt += int64(tx.Seat.Price)
		itemDetail := midtrans.ItemDetails{
			ID:       strconv.FormatUint(uint64(tx.Seat.SeatId), 10),
			Price:    int64(tx.Seat.Price),
			Category: tx.Seat.Category,
			Qty:      1,
			Name:     tx.Seat.Name,
		}
		itemDetails = append(itemDetails, itemDetail)

	}
	var snapRequest snap.Request = snap.Request{ //create snap request data object
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: grossAmt,
		},
		CustomerDetail: &customerDetails,
		Items:          &itemDetails,
	}
	return snapRequest, nil
}

func (s *TransactionService) UpdatePaymentStatus(orderId, vendor, confirmation string) error {
	if result := s.txRepo.UpdatePaymentStatus(orderId, vendor, confirmation); result.Error != nil {
		return result.Error
	}
	return nil
}
