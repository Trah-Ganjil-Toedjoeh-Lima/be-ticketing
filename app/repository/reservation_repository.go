package repository

import (
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r ReservationRepository) GetAllSeats(seats *[]model.Seat) *gorm.DB {
	return r.db.Find(seats)
}
