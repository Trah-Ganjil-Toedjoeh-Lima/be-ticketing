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

func NewSeatService(config *config.AppConfig, seatRepo *repository.SeatRepository, txRepo *repository.TransactionRepository) *SeatService {
	return &SeatService{config: config, seatRepo: seatRepo, txRepo: txRepo}
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
	if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil { //get requested seat
		return result.Error
	}
	//start validation logic
	if seat.Status == "available" { //check from seat table
		return nil
	} else { //if seat table not convincing => check form tx table
		var tx model.Transaction
		if result := s.txRepo.GetBySeat(&tx, seatId).Last(&tx); result.Error != nil { //get the newest transaction data for this seat from tx table. Check if the query returns an error
			return result.Error
		} else if result.RowsAffected < 1 { //double-check the seat status, maybe the cause of  unavailableness is because of 'ghost' reservation
			//if there are no seat data in the tx table, it means that it`s only booked by someone and then did not proceed to the transaction process
			//this case can be caused by irresponsible user that left their reservation but not complete the transaction
			return nil
		}
		if time.Now().After(tx.UpdatedAt.Add(s.config.TransactionMinute)) { //kalo data kursi ada di tabel transaction => cek updated_at. If seat update_at + 15 < time => return nil
			//kalo transaksi sebelumnya "ngambang" maka boleh lanjut
			//transaksi ngambang pada kasus ini disebabkan oleh user yang tidak menyelesaikan/kelamaan dalam proses transaksi
			return nil
		}
		if tx.UserId == userId { // kalo tx sebelumnya gak "ngambang", asalkan yang pesen usernya sama, lanjut
			return nil
		}
		return errors.New("kursi sudah ada yang nge-booking") //kalo gagal melewati constraint diatas, berarti sedang/sudah di cim orang lain
	}
}
