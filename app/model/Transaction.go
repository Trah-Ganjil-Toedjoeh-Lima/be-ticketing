package model

import (
	"gorm.io/gorm"
	"time"
)

type Transaction struct {
	TransactionId int64 `gorm:"primaryKey"`
	MidtransTxId  int64 `gorm:"unique"`
	UserId        int64 `gorm:"not null"`
	SeatId        int   `gorm:"not null"`
	User          User
	Seat          Seat
	Price         string
	Vendor        string
	Confirmation  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}
