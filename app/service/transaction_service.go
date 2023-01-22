package service

import (
	//"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	//"github.com/google/uuid"
)

type TrsansactionService struct {
	txRepo *repository.TransactionRepository
	userRepo *repository.UserRepository

}

func NewTrsansactionService(txRepo *repository.TransactionRepository) *TrsansactionService {
	return &TrsansactionService{txRepo: txRepo}
}

func (s *TrsansactionService) CreateTx(userId uint64, seatIds []uint) {
/*	txId := uuid.New().String()
	for _, seatId := range seatIds {

		newTx := model.Transaction{
			MidtransTxId: txId,
			UserId: userId,
			SeatId: seatId,
			User:,
			Seat: ,
			Price: ,
			Vendor: ,
			Confirmation: ,

		}*/
	}

}
