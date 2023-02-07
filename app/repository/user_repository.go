package repository

import (
	"github.com/frchandra/ticketing-gmcgo/app/model"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *util.LogUtil
}

func NewUserRepository(db *gorm.DB, log *util.LogUtil) *UserRepository {
	return &UserRepository{db: db, log: log}
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
	result := ur.db.Where("name = ? AND email = ? AND phone = ?", user.Name, user.Email, user.Phone).Find(user)
	if result.Error != nil {
		ur.log.BasicLog(result.Error, "UserRepository@FindOne")
	}
	return result
}

func (ur *UserRepository) InsertOne(user *model.User) *gorm.DB {
	result := ur.db.Create(user)
	if result.Error != nil {
		ur.log.BasicLog(result.Error, "UserRepository@InsertOne")
	}
	return result
}

func (ur *UserRepository) GetByPairs(userInput, userOut *model.User) *gorm.DB {
	var result *gorm.DB
	if result = ur.db.Model(model.User{}).Where("name = ?", userInput.Name).Take(userOut); result.Error == nil {
		return result
	}
	if result = ur.db.Model(model.User{}).Where("email = ?", userInput.Email).Take(userOut); result.Error == nil {
		return result
	}
	ur.log.BasicLog(result.Error, "UserRepository@GetByPairs")
	return result
}

func (ur *UserRepository) GetById(userId uint64, userOut *model.User) *gorm.DB {
	result := ur.db.First(userOut, userId)
	if result.Error != nil {
		ur.log.BasicLog(result.Error, "UserRepository@GetById")
	}
	return result

}
