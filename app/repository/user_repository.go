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

func (this *UserRepository) GetByPairs(user *model.User) *gorm.DB {
	var oldUser model.User
	var result *gorm.DB
	if result = this.db.Model(model.User{}).Where("name = ?", user.Name).Take(&oldUser); result.Error == nil {
		return result
	}
	if result = this.db.Model(model.User{}).Where("email = ?", user.Email).Take(&oldUser); result.Error == nil {
		return result
	}
	return result

}
