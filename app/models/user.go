package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId    int `gorm:"primaryKey"`
	Name      string
	Email     string
	phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
