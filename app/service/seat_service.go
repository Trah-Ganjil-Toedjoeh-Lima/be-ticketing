package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/config"
	"time"
)

type SeatService struct {
	config   *config.AppConfig
	seatRepo *repository.SeatRepository
	txRepo   *repository.TransactionRepository
}

func NewSeatService(seatRepo *repository.SeatRepository, txRepo *repository.TransactionRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo, txRepo: txRepo}
}

func (s *SeatService) GetAllSeats() ([]model.Seat, error) {
	var seats []model.Seat
	if err := s.seatRepo.GetAllSeats(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *SeatService) UpdateStatus(seatId uint, status string) error {
	if result := s.seatRepo.UpdateStatus(seatId, status); result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *SeatService) IsOwned(seatId uint, userId uint64) error {
	var seat model.Seat
	//get requested seat
	if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil {
		return result.Error
	}
	//start validation logic
	if seat.Status == "available" { //check from seat table
		return nil
	} else { //if seat table not convincing => check form tx table
		//get the newest transaction data for this seat from tx table
		var tx model.Transaction
		//double-check the seat status, maybe the cause of not availableness is because of 'ghost' reservation
		if result := s.txRepo.GetBySeat(&tx, seatId).Last(&tx); result.Error != nil { //check if the query returns an error
			return result.Error
		} else if result.RowsAffected < 1 { //if there are no seat data in the tx table, it means that it`s only booked by someone and then did not proceed to the transaction process
			return nil //this case can be caused by irresponsible user that left their reservation but not complete the transaction
		}
		//kalo data kursi ada di tabel transaction => cek updated_at. If seat update_at + 15 < time => return nil
		if time.Now().After(tx.UpdatedAt.Add(s.config.TransactionMinute)) {
			//kalo transaksi sebelumnya "ngambang" maka boleh lanjut
			return nil //transaksi ngambang pada kasus ini disebabkan oleh user yang tidak menyelesaikan/kelamaan dalam proses transaksi
		}
		// kalo tx sebelumnya gak "ngambang", asalkan yang pesen usernya sama, lanjut
		if tx.UserId == userId {
			return nil
		}
		//kalo gagal melewati constraint diatas, berarti sedang/sudah di cim orang lain
		return errors.New("kursi sudah ada yang nge-booking")
	}
}
