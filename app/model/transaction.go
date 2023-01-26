package model

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	TransactionId uint64 `gorm:"primaryKey"`
	OrderId       string `gorm:"not null"`
	UserId        uint64 `gorm:"not null"`
	SeatId        uint   `gorm:"not null"`
	User          User
	Seat          Seat
	Vendor        string
	Confirmation  string
	CreatedAt     time.Time      `json:"-"`c

	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}
