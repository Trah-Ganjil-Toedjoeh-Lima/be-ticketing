package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/util"
	"github.com/frchandra/gmcgo/config"
	"strconv"
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

func (s *ReservationService) CheckUserSeatCount(seatIds []uint, userId uint64) error {
	//ambil data transaksi user yang sudah tercatat
	prevTransaction, err := s.txService.GetTxDetailsByUser(userId)
	if err != nil {
		return err
	}
	//ambil data seatId nya saja
	var prevTxSeatIds []uint
	for _, tx := range prevTransaction {
		prevTxSeatIds = append(prevTxSeatIds, tx.SeatId)
	}
	//ambil perbedaan seatId sesudah dan sebelum
	diff := util.ElementDifference(seatIds, prevTxSeatIds)
	//jika jumlah kursi sebelumnya + jumlah kursi pesanan yang belum ada di transaksi sebelumnya > 5 return error
	if totalSeat := len(diff) + len(prevTransaction); totalSeat > 5 {
		return errors.New("user telah memesan " + strconv.Itoa(len(prevTransaction)) + " kursi, tidak bisa memesan " + strconv.Itoa(len(seatIds)) + " kursi lagi")
	}
	return nil
}
