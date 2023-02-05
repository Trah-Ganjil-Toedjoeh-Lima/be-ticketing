package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId      uint64         `gorm:"primaryKey"`
	Name        string         `gorm:"not null"`
	Email       string         `gorm:"not null"`
	Phone       string         `gorm:"not null"`
	Transaction []Transaction  `gorm:"foreignKey:UserId;references:UserId"json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
