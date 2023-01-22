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

func (ur *UserRepository) GetOrInsertOne(user *model.User) *gorm.DB {
	result := ur.FindOne(user)
	if result.RowsAffected < 1 {
		return ur.InsertOne(user)
	} else {
		return result
	}

}

func (ur *UserRepository) FindOne(user *model.User) *gorm.DB {
	return ur.db.Where("name = ? AND email = ? AND phone = ?", user.Name, user.Email, user.Phone).Find(user)
}

func (ur *UserRepository) InsertOne(user *model.User) *gorm.DB {
	result := ur.db.Create(user)
	return result
}

func (ur *UserRepository) GetByPairs(userInput, userOut *model.User) *gorm.DB {
	var result *gorm.DB //TODO: this not need "Model(model.User{})....
	if result = ur.db.Model(model.User{}).Where("name = ?", userInput.Name).Take(userOut); result.Error == nil {
		return result
	}
	if result = ur.db.Model(model.User{}).Where("email = ?", userInput.Email).Take(userOut); result.Error == nil {
		return result
	}
	return result
}

func (ur *UserRepository) GetById(userId uint64, userOut *model.User) *gorm.DB {
	result := ur.db.First(userOut, userId)
	return result

}
