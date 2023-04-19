package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId      uint64 `gorm:"primaryKey"`
	Name        string
	Email       string `gorm:"not null;uniqueIndex"`
	Phone       string
	TotpSecret  string         `json:"-"`
	Transaction []Transaction  `gorm:"foreignKey:UserId;references:UserId"json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
