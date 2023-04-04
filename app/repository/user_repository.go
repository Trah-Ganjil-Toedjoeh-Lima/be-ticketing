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

func (u *UserRepository) GetOrInsertOne(user *model.User) *gorm.DB {
	result := u.FindOne(user)
	if result.RowsAffected < 1 {
		return u.InsertOne(user)
	} else {
		return result
	}

}

func (u *UserRepository) FindOne(user *model.User) *gorm.DB {
	result := u.db.Where("name = ? AND email = ? AND phone = ?", user.Name, user.Email, user.Phone).Find(user)
	if result.Error != nil {
		u.log.BasicLog(result.Error, "UserRepository@FindOne")
	}
	return result
}

func (u *UserRepository) InsertOne(user *model.User) *gorm.DB {
	result := u.db.Create(user)
	if result.Error != nil {
		u.log.BasicLog(result.Error, "UserRepository@InsertOne")
	}
	return result
}

func (u *UserRepository) GetByPairs(userInput, userResult *model.User) *gorm.DB {
	var result *gorm.DB
	if result = u.db.Model(model.User{}).Where("name = ?", userInput.Name).Take(userResult); result.Error == nil {
		return result
	}
	if result = u.db.Model(model.User{}).Where("email = ?", userInput.Email).Take(userResult); result.Error == nil {
		return result
	}
	u.log.BasicLog(result.Error, "UserRepository@GetByPairs")
	return result
}

func (u *UserRepository) GetByEmail(email string, userResult *model.User) *gorm.DB {
	result := u.db.Where(model.User{}).Where("email = ?", email).First(userResult)
	if result.Error != nil {
		u.log.BasicLog(result.Error, "-")
	}
	return result
}

func (u *UserRepository) GetById(userId uint64, userResult *model.User) *gorm.DB {
	result := u.db.First(userResult, userId)
	if result.Error != nil {
		u.log.BasicLog(result.Error, "UserRepository@GetById")
	}
	return result

}
