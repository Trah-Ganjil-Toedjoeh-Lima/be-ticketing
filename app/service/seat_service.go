package service

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/repository"
	"github.com/frchandra/ticketing-gmcgo/config"
	"gorm.io/gorm"
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
	if err := s.seatRepo.GetAll(&seats).Error; err != nil {
		return nil, errors.New("database operation error")
	}
	return seats, nil
}

func (s *SeatService) UpdatePostSaleStatus(link, status string) error {
	if result := s.seatRepo.UpdatePostSaleStatus(link, status); result.Error != nil {
		return errors.New("database operation error")
	}
	return nil
}

func (s *SeatService) UpdateStatus(seatId uint, status string) error {
	if result := s.seatRepo.UpdateStatus(seatId, status); result.Error != nil {
		return errors.New("database operation error")
	}
	return nil
}

func (s *SeatService) UpdateStatusTxn(txn *gorm.DB, seatId uint, status string) error {
	if result := s.seatRepo.UpdateStatusTxn(txn, seatId, status); result.Error != nil {
		return errors.New("database operation error")
	}
	return nil
}

func (s *SeatService) IsOwnedTxn(txn *gorm.DB, seatId uint, userId uint64) error {
	var seat model.Seat
	if result := s.seatRepo.GetByIdTxn(txn, &seat, seatId); result.Error != nil { //get requested seat
		return errors.New("database operation error")
	}
	//start validation logic
	if seat.Status == "available" { //check from seat table
		return nil
	} else { //if seat table is not convincing => check form tx table
		if seat.Status == "not_for_sale" {
			return errors.New("this seat is not for sale (nice try wkwkwk)")
		}

		var transaction model.Transaction
		if result := s.txRepo.GetBySeatTxn(txn, &transaction, seatId).Last(&transaction); result.Error != nil { //get the newest transaction data for this seat from transaction table. Check if the query returns an error
			return errors.New("database operation error")
		} else if result.RowsAffected < 1 { //double-check the seat status, maybe the cause of  unavailableness is because of 'ghost' reservation
			//if there are no seat data in the transaction table, it means that it`s only booked by someone and then did not proceed to the transaction process
			//this case can be caused by irresponsible user that left their reservation but not complete the transaction
			return nil
		}

		if transaction.Confirmation == "settlement" { //if this transaction is already settled it mean that this seat is unavailable
			return errors.New("this seat is not available")
		} else if time.Now().After(transaction.UpdatedAt.Add(s.config.TransactionMinute)) { //if seat not settled, then continue to check. Cek data kursi ada di tabel transaction => cek updated_at. If seat update_at + 15 < time => return nil
			//kalo transaksi sebelumnya "ngambang" maka boleh lanjut
			//transaksi ngambang pada kasus ini disebabkan oleh user yang tidak menyelesaikan/kelamaan dalam proses transaksi
			return nil
		}

		if transaction.UserId == userId { // kalo tx sebelumnya gak "ngambang", asalkan yang pesen usernya sama, lanjut
			return nil
		}

		return errors.New("this seat is already booked") //kalo gagal melewati constraint diatas, berarti sedang/sudah di cim orang lain
	}
}
