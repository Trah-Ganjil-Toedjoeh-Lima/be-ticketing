package model

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	TransactionId uint64 `gorm:"primaryKey"`
	MidtransTxId  string `gorm:"unique"`
	UserId        uint64 `gorm:"not null"`
	SeatId        uint   `gorm:"not null"`
	User          User
	Seat          Seat
	Price         string
	Vendor        string
	Confirmation  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
