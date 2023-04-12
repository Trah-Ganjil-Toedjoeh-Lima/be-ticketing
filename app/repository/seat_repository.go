package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) UpdatePostSaleStatus(link, status string) *gorm.DB {
	result := r.db.Model(&model.Seat{}).Where("link = ?", link).Update("post_sale_status", status)
	return result
}

func (r *SeatRepository) UpdateStatus(seatId uint, status string) *gorm.DB {
	result := r.db.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
	return result
}

func (r *SeatRepository) UpdateStatusTxn(txn *gorm.DB, seatId uint, status string) *gorm.DB {
	result := txn.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
	return result
}

func (r *SeatRepository) GetAll(seats *[]model.Seat) *gorm.DB {
	result := r.db.Find(seats)
	return result
}

func (r *SeatRepository) GetByIdTxn(txn *gorm.DB, seat *model.Seat, id uint) *gorm.DB {
	result := txn.First(seat, id)
	return result
}
