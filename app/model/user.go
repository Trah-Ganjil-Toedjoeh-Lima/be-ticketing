package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserId      uint64 `gorm:"primaryKey"`
	Name        string
	Email       string `gorm:"not null"`
	Phone       string
	TotpSecret1 string         `json:"-"`
	TotpSecret2 string         `json:"-"`
	Transaction []Transaction  `gorm:"foreignKey:UserId;references:UserId"json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
