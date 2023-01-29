package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/config"
	"strconv"
	"time"
)

type ReservationService struct {
	seatRepo  *repository.SeatRepository
	txRepo    *repository.TransactionRepository
	txService *TransactionService
	config    *config.AppConfig
}

func NewReservationService(seatRepo *repository.SeatRepository, txRepo *repository.TransactionRepository, txService *TransactionService, config *config.AppConfig) *ReservationService {
	return &ReservationService{seatRepo: seatRepo, txRepo: txRepo, txService: txService, config: config}
}

func (s *ReservationService) GetAllSeats() ([]model.Seat, error) {
	seats := []model.Seat{}
	if err := s.seatRepo.GetAllSeats(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *ReservationService) IsOwned(seatId uint, userId uint64) error {
	var seat model.Seat
	//get requested seat
	if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil {
		return result.Error
	}
	//start validation logic
	if seat.Status == "available" {
		return nil
	} else {
		//ambil tx terbaru untuk kursi ini
		var tx model.Transaction
		//cek lebih detail status kursi, barangkali not availablenya gara-gara transaksi ngambang
		if result := s.txRepo.GetBySeat(&tx, seatId).Last(&tx); result.Error != nil {
			return result.Error
		} else if result.RowsAffected < 1 { //kalo gak ada berarti aman, lanjut
			return nil
		}
		//kalo ada cek updated_at: if seat update_at + 15 < time now OR time_now - updated_at > 15m-> return nil
		if time.Now().After(tx.UpdatedAt.Add(s.config.TransactionMinute)) {
			//kalo transaksi sebelumnya "ngambang" maka boleh lanjut
			return nil
		}
		// kalo tx sebelumnya gak "ngambang", asalkan yang pesen usernya sama, lanjut
		if tx.UserId == userId {
			return nil
		}
		//kalo gagal melewati constraint diatas, berarti sedang/sudah di cim orang lain
		return errors.New("kursi sudah ada yang nge-booking")
	}
}

func (s *ReservationService) CheckUserSeatCount(seatIds []uint, userId uint64) error {
	//ambil data transaksi user yang sudah tercatat
	prevTransaction, _ := s.txService.GetTxDetailsByUser(userId)
	//ambil data seatId nya saja
	var prevTxSeatIds []uint
	for _, tx := range prevTransaction {
		prevTxSeatIds = append(prevTxSeatIds, tx.SeatId)
	}
	//ambil perbedaan seatId sesudah dan sebelum
	diff := s.difference(seatIds, prevTxSeatIds)
	//jika jumlah kursi sebelumnya + jumlah kursi pesanan yang belum ada di transaksi sebelumnya > 5 return error
	if totalSeat := len(diff) + len(prevTransaction); totalSeat > 5 {
		return errors.New("user telah memesan " + strconv.Itoa(len(prevTransaction)) + " kursi, tidak bisa memesan " + strconv.Itoa(len(seatIds)) + " kursi lagi")
	}
	return nil
}

func (s *ReservationService) difference(after, before []uint) []uint {
	mb := make(map[uint]struct{}, len(before))
	for _, x := range before {
		mb[x] = struct{}{}
	}
	var diff []uint
	for _, x := range after {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
