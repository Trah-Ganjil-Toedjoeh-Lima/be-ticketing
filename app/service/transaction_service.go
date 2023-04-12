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
			return errors.New("database operation error")
		}
	}
	return nil
}

func (s *TransactionService) GetAllWithDetails() ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetAllWithDetails(&transactions); result.Error != nil {
		return transactions, errors.New("database operation error")
	}
	return transactions, nil
}

func (s *TransactionService) GetDetailsByLink(link string) (model.Transaction, error) {
	var transaction model.Transaction
	if result := s.txRepo.GetDetailsByLink(&transaction, link); result.Error != nil {
		return transaction, errors.New("database operation error")
	}
	return transaction, nil
}

func (s *TransactionService) GetDetailsByUser(userId uint64) ([]model.Transaction, error) {
	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUser(&transactions, userId); result.Error != nil {
		return transactions, errors.New("database operation error")
	}
	transactions = s.CleanUpGhostTransaction(transactions)
	return transactions, nil
}

func (s *TransactionService) GetDetailsByUserConfirmation(userId uint64, confirmation string) ([]model.Transaction, error) {
	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUserConfirmation(&transactions, userId, confirmation); result.Error != nil {
		return transactions, errors.New("database operation error")
	}
	transactions = s.CleanUpGhostTransaction(transactions)
	return transactions, nil
}

func (s *TransactionService) GetByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetByOrder(&transactions, orderId); result.Error != nil {
		return transactions, errors.New("database operation error")
	}
	return transactions, nil
}

func (s *TransactionService) GetDetailsByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetDetailsByOrder(&transactions, orderId); result.Error != nil {
		return transactions, errors.New("database operation error")
	}
	return transactions, nil
}

func (s *TransactionService) CleanUpGhostTransaction(transactions []model.Transaction) []model.Transaction {
	var newTransaction []model.Transaction //cek apakah ada transaksi ngambang, jika ada buang dari slice dan update db
	for _, tx := range transactions {
		if time.Now().After(tx.CreatedAt.Add(s.config.TransactionMinute)) && tx.Confirmation != "settlement" { //if tx created_at + 15 < time now  => berarti transaction ngambang
			//update database
			s.txRepo.UpdateUserPaymentStatus(tx.UserId, "", "not_continued")
			s.txRepo.SoftDeleteBySeatUser(tx.SeatId, tx.UserId)
		} else {
			newTransaction = append(newTransaction, tx)
		}
	}
	//transaksi bersih
	return newTransaction
}

func (s *TransactionService) SeatsBelongsToUser(userId uint64) ([]model.Seat, error) {
	var seats []model.Seat

	var transactions []model.Transaction //get user's transaction
	if result := s.txRepo.GetDetailsByUser(&transactions, userId); result.Error != nil {
		return seats, errors.New("database operation error")
	}
	if transactions = s.CleanUpGhostTransaction(transactions); len(transactions) < 1 {
		return seats, errors.New("this user doesen`t have any transaction")
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
	orderId := uuid.New().String()               //create order_id for the new midtrans transaction
	s.txRepo.UpdateUserOrderId(userId, orderId)  //update order_id of this transaction in the database
	customerDetails := midtrans.CustomerDetails{ //populate the midtrans request with the customer detail
		FName: txDetails[0].User.Name,
		LName: "",
		Email: txDetails[0].User.Email,
		Phone: txDetails[0].User.Phone,
	}
	var grossAmt int64 //populate the item detail
	var itemDetails []midtrans.ItemDetails
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
		return errors.New("database operation error")
	}
	return nil
}
