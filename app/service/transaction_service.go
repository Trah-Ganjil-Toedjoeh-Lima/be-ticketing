package service

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/google/uuid"
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
	var user model.User
	if result := s.userRepo.GetById(userId, &user); result.Error != nil {
		return result.Error
	}

	for _, seatId := range seatIds {
		var seat model.Seat
		if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil {
			return result.Error
		}
		newTx := model.Transaction{
			MidtransTxId: txId,
			UserId:       userId,
			SeatId:       seatId,
			User:         user,
			Seat:         seat,
			Vendor:       "#",
			Confirmation: "#",
		}

		if result := s.txRepo.InsertOne(&newTx); result.Error != nil {
			return result.Error
		}
	}
	return nil
}
