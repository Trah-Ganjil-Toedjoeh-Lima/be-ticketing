package repository

import (
	"github.com/frchandra/gmcgo/app/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (this *UserRepository) InsertOne(user *model.User) *gorm.DB {
	result := this.db.Create(user)
	return result
}
