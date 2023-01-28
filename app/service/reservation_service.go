package service

import (
	"errors"
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"strconv"
)

type ReservationService struct {
	seatRepo  *repository.SeatRepository
	txRepo    *repository.TransactionRepository
	txService *TransactionService
}

func NewReservationService(seatRepo *repository.SeatRepository, txRepo *repository.TransactionRepository, txService *TransactionService) *ReservationService {
	return &ReservationService{seatRepo: seatRepo, txRepo: txRepo, txService: txService}
}

func (s *ReservationService) GetAllSeats() ([]model.Seat, error) {
	seats := []model.Seat{}
	if err := s.seatRepo.GetAllSeats(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *ReservationService) IsOwned(seatId uint, userId uint64) error {
	var seat model.Seat //TODO: think of all edge scenarios
	//get requested seat
	if result := s.seatRepo.GetSeatById(&seat, seatId); result.Error != nil {
		return result.Error
	}
	//if kursi masih kosong -> return ok
	if seat.Status == "#" {
		return nil
	} else {
		var tx model.Transaction
		result := s.txRepo.GetLastTxBySeatIdUserId(&tx, seat.SeatId, userId)
		//if this seat has been reserved by this user //if sudah ada yang ngisi tapi dirinya sendiri -> return ok (asalkan transaksi belum berjalan/selesai)
		if result.RowsAffected == 1 {
			//if transaction has been done or during the transaction
			if tx.Confirmation == "settlement" || tx.Confirmation == "pending" {
				return errors.New("seat is pending or payed")
			} else {
				return nil
			}

		} else { //the seat is booked but not by this user
			return errors.New("kursi sudah ada yang nge-booking")
		}
	}
}

func (s *ReservationService) CheckUserSeatCount(seatIds []uint, userId uint64) error {
	//ambil data transaksi user yang sudah tercatat
	prevTransaction, _ := s.txService.GeTxDetailsByUser(userId)
	//ambil data seatId nya saja
	var prevTxSeatIds []uint
	for _, tx := range prevTransaction {
		prevTxSeatIds = append(prevTxSeatIds, tx.SeatId)
	}
	//ambil perbedaan seatId sesudah dan sebelum
	diff := difference(seatIds, prevTxSeatIds)
	//jika jumlah kursi sebelumnya + jumlah kursi pesanan yang belum ada di transaksi sebelumnya > 5 return error
	if totalSeat := len(diff) + len(prevTransaction); totalSeat > 5 {
		return errors.New("user telah memesan " + strconv.Itoa(len(prevTransaction)) + " kursi, tidak bisa memesan " + strconv.Itoa(len(seatIds)) + " kursi lagi")
	}
	return nil
}

func difference(after, before []uint) []uint {
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
