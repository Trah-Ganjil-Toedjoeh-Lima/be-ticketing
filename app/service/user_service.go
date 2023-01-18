package service

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/util/token"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (this *UserService) InsertOne(user *model.User) (int64, error) {
	hashedCred, err := bcrypt.GenerateFromPassword([]byte(user.Phone), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Phone = string(hashedCred)
	result := this.userRepository.InsertOne(user)
	if result.Error != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func (this *UserService) verifyCredentials(cred, hashedCred string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedCred), []byte(cred))
}

func (this *UserService) ValidateLogin(userInput *model.User) (string, error) {
	var oldUser model.User
	result := this.userRepository.GetByPairs(userInput)
	if result.Error != nil {
		return "", result.Error
	}
	result.First(&oldUser)
	err := this.verifyCredentials(userInput.Phone, oldUser.Phone)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	token, err := token.GenerateToken(oldUser.UserId)
	if err != nil {
		return "", err
	}
	return token, nil

}
