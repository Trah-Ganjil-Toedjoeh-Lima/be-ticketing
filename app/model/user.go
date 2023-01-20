package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId      uint64        `gorm:"primaryKey"`
	Name        string        `gorm:"not null"`
	Email       string        `gorm:"not null"`
	Phone       string        `gorm:"not null"`
	Transaction []Transaction `gorm:"foreignKey:UserId;references:UserId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
