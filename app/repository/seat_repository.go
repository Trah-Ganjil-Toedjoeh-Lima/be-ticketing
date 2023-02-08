package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db  *gorm.DB
	log *util.LogUtil
}

func NewSeatRepository(db *gorm.DB, log *util.LogUtil) *SeatRepository {
	return &SeatRepository{db: db, log: log}
}

func (r *SeatRepository) UpdateStatus(seatId uint, status string) *gorm.DB {
	result := r.db.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
	if result.Error != nil {
		r.log.BasicLog(result.Error, "SeatRepository@UpdateStatusTxn")
	}
	return result
}

func (r *SeatRepository) UpdateStatusTxn(txn *gorm.DB, seatId uint, status string) *gorm.DB {
	result := txn.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
	if result.Error != nil {
		r.log.BasicLog(result.Error, "SeatRepository@UpdateStatusTxn")
	}
	return result
}

func (r *SeatRepository) GetAll(seats *[]model.Seat) *gorm.DB {
	result := r.db.Find(seats)
	if result.Error != nil {
		r.log.BasicLog(result.Error, "SeatRepository@GetAll")
	}
	return result
}

func (r *SeatRepository) GetByIdTxn(txn *gorm.DB, seat *model.Seat, id uint) *gorm.DB {
	result := txn.First(seat, id)
	if result.Error != nil {
		r.log.BasicLog(result.Error, "SeatRepository@GetByIdTxn")
	}
	return result
}
