package model

import (
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	SeatId      uint          `gorm:"primaryKey"`
	Name        string        `gorm:"unique;not null"`
	Price       uint          `gorm:"not null"`
	Link        string        `gorm:"not null"`
	Status      string        `gorm:"not null"`
	Transaction []Transaction `gorm:"foreignKey:SeatId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
