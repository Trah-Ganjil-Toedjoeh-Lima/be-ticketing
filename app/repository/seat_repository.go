package repository

import (
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) GetSeatById(seat *model.Seat, seatId uint) *gorm.DB {
	return r.db.Where(seat, seatId)
}

func (r *SeatRepository) UpdateStatus(seatId uint, status string) *gorm.DB {
	return r.db.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
}
