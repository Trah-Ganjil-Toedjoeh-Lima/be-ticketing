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

func (r *SeatRepository) Atomic(
	fn func(newSeatRepo *SeatRepository) *gorm.DB,
) (result *gorm.DB) {
	txDb := r.db.Begin()
	defer func() {
		if result.Error != nil {
			txDb.Rollback()
		} else {
			txDb.Commit()
		}
	}()
	newSeatRepo := NewSeatRepository(txDb)
	result = fn(newSeatRepo)
	return result

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
