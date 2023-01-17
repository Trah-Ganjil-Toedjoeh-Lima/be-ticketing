package model

import (
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	SeatId    int    `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Price     int    `gorm:"not null"`
	Link      string `gorm:"not null"`
	Status    int8   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
