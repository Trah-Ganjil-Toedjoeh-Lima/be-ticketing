package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/config"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"strconv"
	"time"
)

type TransactionService struct {
	txRepo   *repository.TransactionRepository
	userRepo *repository.UserRepository
	seatRepo *repository.SeatRepository
	config   *config.AppConfig
}

func NewTransactionService(txRepo *repository.TransactionRepository, userRepo *repository.UserRepository, seatRepo *repository.SeatRepository, config *config.AppConfig) *TransactionService {
	return &TransactionService{txRepo: txRepo, userRepo: userRepo, seatRepo: seatRepo, config: config}
}

func (s *TransactionService) CreateTx(userId uint64, seatIds []uint) error {
	for _, seatId := range seatIds {
		//create tx for each seat
		newTx := model.Transaction{
			OrderId:      "",
			UserId:       userId,
			SeatId:       seatId,
			Vendor:       "no_vendor",
			Confirmation: "reserved",
		}
		//delete previous failed reservation
		s.txRepo.SoftDeleteBySeatUser(seatId, userId)
		//save transaction
		if result := s.txRepo.InsertOne(&newTx); result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (s *TransactionService) GetTxDetailsByUser(userId uint64) ([]model.Transaction, error) {
	//get user's transaction
	var transactions []model.Transaction
	if result := s.txRepo.GetByUser(&transactions, userId); result.Error != nil {
		return transactions, result.Error
	}
	transactions = s.CleanUpDeadTransaction(transactions)
	return transactions, nil
}

func (s *TransactionService) GetTxByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetByOrder(&transactions, orderId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) GetTxDetailsByOrder(orderId string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	if result := s.txRepo.GetDetailsByOrder(&transactions, orderId); result.Error != nil {
		return transactions, result.Error
	}
	return transactions, nil
}

func (s *TransactionService) CleanUpDeadTransaction(transactions []model.Transaction) []model.Transaction {
	var newTransaction []model.Transaction
	//cek apakah ada transaksi ngambang, jika ada buang dari slice dan update db
	for _, tx := range transactions {
		//if tx update_at + 15 < time now  => berarti transaction ngambang
		if time.Now().After(tx.UpdatedAt.Add(s.config.TransactionMinute)) {
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

func (s *TransactionService) IsSeatsBelongsToUser(userId uint64) ([]model.Seat, error) {
	var seats []model.Seat
	//get user's transaction
	var transactions []model.Transaction
	if result := s.txRepo.GetDetailsByUser(&transactions, userId); result.Error != nil {
		return seats, result.Error
	}
	if transactions = s.CleanUpDeadTransaction(transactions); len(transactions) < 1 {
		return seats, errors.New("this user doesen`t have any transaction")
	}
	for _, tx := range transactions {
		var seat model.Seat
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

func (s *TransactionService) PrepareTransactionData(userId uint64) (snap.Request, error) {
	//get user's transaction
	var txDetails []model.Transaction
	s.txRepo.GetDetailsByUser(&txDetails, userId)
	//clean up tx
	if txDetails = s.CleanUpDeadTransaction(txDetails); len(txDetails) < 1 {
		return snap.Request{}, errors.New("cannot find any transaction for this user")
	}
	//create order_id
	orderId := uuid.New().String()
	//update order_id
	s.txRepo.UpdateUserOrderId(userId, orderId)
	//create customer detail
	customerDetails := midtrans.CustomerDetails{
		FName: txDetails[0].User.Name,
		LName: "",
		Email: txDetails[0].User.Email,
		Phone: txDetails[0].User.Phone,
	}
	//create item detail
	var grossAmt int64
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
	//create snap request
	var snapRequest snap.Request = snap.Request{
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
