package model

import (
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	SeatId         uint   `gorm:"primaryKey"`
	Name           string `gorm:"unique;not null"`
	Row            string `gorm:"not null"`
	Column         uint   `gorm:"not null"`
	Price          uint   `gorm:"not null"`
	Link           string `gorm:"not null"`
	Status         string `gorm:"not null"`
	PostSaleStatus string
	Transaction    []Transaction  `gorm:"foreignKey:SeatId"json:"-"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}
