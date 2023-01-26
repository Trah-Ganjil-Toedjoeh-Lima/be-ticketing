package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
)

type ReservationService struct {
	resRepo *repository.ReservationRepository
	txRepo  *repository.TransactionRepository
}

func NewReservationService(resRepo *repository.ReservationRepository, txRepo *repository.TransactionRepository) *ReservationService {
	return &ReservationService{resRepo: resRepo, txRepo: txRepo}
}

func (s *ReservationService) GetAllSeats() ([]model.Seat, error) {
	seats := []model.Seat{}
	if err := s.resRepo.GetAllSeats(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *ReservationService) IsOwned(seatId uint, userId uint64) error {
	//if kursi masih kosong -> return ok
	//else
	//if sudah ada yang ngisi tapi dirinya sendiri -> return ok
	//return error
	var seat model.Seat //TODO: think of all edge scenarios
	if result := s.resRepo.GetSeatById(&seat, seatId); result.Error != nil {
		return result.Error
	}
	if seat.Status == "#" {
		return nil
	} else {
		var tx model.Transaction
		result := s.txRepo.GetLastTxBySeatIdUserId(&tx, seat.SeatId, userId)
		if result.RowsAffected == 1 {
			if tx.Confirmation != "payed" || tx.Confirmation != "pending" {
				return nil
			}
		}
		return errors.New("kursi sudah ada yang nge-booking")
	}
}
