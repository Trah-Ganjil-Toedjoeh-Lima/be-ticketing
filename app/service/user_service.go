package service

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (this *UserService) InsertOne(user *model.User) (int64, error) {
	hashedPhone, err := bcrypt.GenerateFromPassword([]byte(user.Phone), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Phone = string(hashedPhone)
	result := this.userRepository.InsertOne(user)
	if result.Error != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
