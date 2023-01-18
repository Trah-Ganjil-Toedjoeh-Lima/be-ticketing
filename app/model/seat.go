package model

import (
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	SeatId    uint   `gorm:"primaryKey"`
	Name      string `gorm:"unique;not null"`
	Price     uint   `gorm:"not null"`
	Link      string `gorm:"not null"`
	Status    uint8  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
