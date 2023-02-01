package repository

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SeatRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewSeatRepository(db *gorm.DB, logger *logrus.Logger) *SeatRepository {
	return &SeatRepository{db: db, logger: logger}
}

func (r *SeatRepository) UpdateStatus(seatId uint, status string) *gorm.DB {
	return r.db.Model(&model.Seat{}).Where("seat_id = ?", seatId).Update("status", status)
}

func (r *SeatRepository) GetAllSeats(seats *[]model.Seat) *gorm.DB {
	return r.db.Find(seats)
}

func (r *SeatRepository) GetSeatById(seat *model.Seat, id uint) *gorm.DB {
	return r.db.First(seat, id)
}
