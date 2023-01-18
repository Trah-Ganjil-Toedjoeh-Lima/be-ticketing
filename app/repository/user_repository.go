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

func (this *UserRepository) GetByPairs(userInput, userOut *model.User) *gorm.DB {
	var result *gorm.DB
	if result = this.db.Model(model.User{}).Where("name = ?", userInput.Name).Take(userOut); result.Error == nil {
		return result
	}
	if result = this.db.Model(model.User{}).Where("email = ?", userInput.Email).Take(userOut); result.Error == nil {
		return result
	}
	return result
}

func (this *UserRepository) GetById(userId uint, userOut *model.User) *gorm.DB {
	result := this.db.First(userOut, userId)
	return result

}
