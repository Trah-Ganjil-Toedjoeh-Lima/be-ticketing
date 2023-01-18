package service

import (
	"github.com/frchandra/gmcgo/app/model"
	"github.com/frchandra/gmcgo/app/repository"
	"github.com/frchandra/gmcgo/app/util"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repository.UserRepository
	tokenUtil      *util.TokenUtil
}

func NewUserService(userRepository *repository.UserRepository, tokenUtil *util.TokenUtil) *UserService {
	return &UserService{
		tokenUtil:      tokenUtil,
		userRepository: userRepository,
	}
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

func (this *UserService) ValidateLogin(userInput *model.User) (*util.TokenDetails, error) {
	var userOut model.User
	var err error
	var tokenDetails *util.TokenDetails
	result := this.userRepository.GetByPairs(userInput, &userOut)
	if result.Error != nil {
		return tokenDetails, result.Error
	}
	err = this.verifyCredentials(userInput.Phone, userOut.Phone)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return tokenDetails, err
	}
	//jwt, err := token.GenerateToken(userOut.UserId)
	tokenDetails, err = this.tokenUtil.CreateToken(userOut.UserId)
	if err != nil {
		return tokenDetails, err
	}
	if err = this.tokenUtil.StoreAuthn(userOut.UserId, tokenDetails); err != nil {
		return tokenDetails, err
	}
	return tokenDetails, nil
}

func (this *UserService) GetById(userId uint64) (model.User, error) {
	var userOut model.User
	result := this.userRepository.GetById(userId, &userOut)
	if result.Error != nil {
		return userOut, result.Error
	}
	return userOut, nil
}
