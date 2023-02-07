package service

import (
	"errors"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"strconv"
)

type ReservationService struct {
	config    *config.AppConfig
	txService *TransactionService
}

func NewReservationService(config *config.AppConfig, txService *TransactionService) *ReservationService {
	return &ReservationService{config: config, txService: txService}
}

func (s *ReservationService) CheckUserSeatCount(seatIds []uint, userId uint64) error {
	prevTransaction, err := s.txService.GetTxDetailsByUser(userId) //ambil data transaksi user yang sudah tercatat
	if err != nil {
		return errors.New("database operation error")
	}

	var prevTxSeatIds []uint //ambil data seatId nya saja
	for _, tx := range prevTransaction {
		prevTxSeatIds = append(prevTxSeatIds, tx.SeatId)
	}

	diff := util.ElementDifference(seatIds, prevTxSeatIds)            //ambil perbedaan seatId sesudah dan sebelum
	if totalSeat := len(diff) + len(prevTransaction); totalSeat > 5 { //jika jumlah kursi sebelumnya + jumlah kursi pesanan yang belum ada di transaksi sebelumnya > 5 return error
		return errors.New("user telah memesan " + strconv.Itoa(len(prevTransaction)) + " kursi, tidak bisa memesan " + strconv.Itoa(len(seatIds)) + " kursi lagi")
	}
	return nil
}
